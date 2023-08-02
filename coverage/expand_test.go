package coverage_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/jsonreference"
	openapispec "github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/coverage"
)

func TestExpand_MediaTransform(t *testing.T) {
	modelName := "Transform"
	modelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"
	model, err := coverage.Expand(modelName, modelSwaggerPath)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("expanded model %s", string(out))
}

func TestExpand_referToChildVariant(t *testing.T) {
	swaggerPath := "./testdata/"
	pathRef := jsonreference.MustCreateRef("test1.json#/paths/~1path1/put")
	swaggerModel, err := coverage.GetModelInfoFromIndexRef(openapispec.Ref{Ref: pathRef}, swaggerPath)
	if err != nil {
		t.Fatal(err)
	}

	if swaggerModel.ModelName != "pet" {
		t.Fatalf("expected modelName pet, got %s", swaggerModel.ModelName)
	}

	model, err := coverage.Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("expanded model %s", string(out))

	if model.Properties == nil {
		t.Fatalf("expected properties not nil")
	}

	if _, ok := (*model.Properties)["odata.type"]; !ok {
		t.Fatalf("expected properties odata.type string")
	}

	if model.Variants == nil {
		t.Fatalf("expected variants not nil")
	}

	if model.Discriminator == nil || *model.Discriminator != "odata.type" {
		t.Fatalf("expected discriminator odata.type")
	}

	if model.VariantType == nil || *model.VariantType != "animal.pet" {
		t.Fatalf("expected variantType animal.pet")
	}

	if _, ok := (*model.Variants)["dog"]; !ok {
		t.Fatalf("expected variants dog not nil")
	}

	if (*model.Variants)["dog"].Properties == nil {
		t.Fatalf("expected variants dog properties not nil")
	}

	if (*model.Variants)["dog"].VariantType == nil || *(*model.Variants)["dog"].VariantType != "animal.pet.dog" {
		t.Fatalf("expected variants dog variantType animal.pet.dog")
	}

	if _, ok := (*(*model.Variants)["dog"].Properties)["odata.type"]; !ok {
		t.Fatalf("expected variants dog properties odata.type string")
	}

	if _, ok := (*(*model.Variants)["dog"].Properties)["name"]; !ok {
		t.Fatalf("expected variants dog properties name")
	}

	if _, ok := (*(*model.Variants)["dog"].Properties)["is_barking"]; !ok {
		t.Fatalf("expected variants dog properties name")
	}

}

// try to expand all PUT and POST models twice, and ensure result is the same
// AZURE_REST_REPO_DIR="/home/test/go/src/github.com/azure/azure-rest-api-specs/specification" TEST_RESULT_FILE="/home/test/"
func TestExpandAll(t *testing.T) {
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
	}
	if strings.HasSuffix(azureRepoDir, "specification") {
		azureRepoDir += "/"
	}
	if !strings.HasSuffix(azureRepoDir, "specification/") {
		t.Fatalf("AZURE_REST_REPO_DIR must specify the specification folder, e.g., AZURE_REST_REPO_DIR=\"/home/test/go/src/github.com/azure/azure-rest-api-specs/specification/\"")
	}

	t.Logf("azure repo dir: %s", azureRepoDir)

	testResultFilePath := os.Getenv("TEST_RESULT_PATH")

	res := testExpandAll(t, azureRepoDir, resultFile(testResultFilePath))
	res2 := testExpandAll(t, azureRepoDir, resultFile(testResultFilePath))
	if res.AllRef != res2.AllRef || res.AllProp != res2.AllProp || res.AvailRef != res2.AvailRef {
		t.Fatalf("result not equal, res1: %+v, res2: %+v", res, res2)
	}
}

func resultFile(testResultFilePath string) string {
	return filepath.Join(testResultFilePath, fmt.Sprintf("res_%s.json", time.Now().Format(time.RFC3339)))
}

type result struct {
	AllRef   int            `json:"all_ref"`
	AllProp  int            `json:"all_prop"`
	AvailRef int            `json:"avail_ref"`
	Paths    map[string]int `json:"paths"`
}

func testExpandAll(t *testing.T, azureRepoDir, testResultPath string) result {
	refList := getRefList(t)
	for _, ref := range refList {
		t.Logf("ref: %s", ref.String())
	}

	// expand concurrently
	res := result{
		AllRef:   len(refList),
		AllProp:  0,
		AvailRef: 0,
		Paths:    make(map[string]int),
	}
	refChan := make(chan jsonreference.Ref)
	lock := sync.RWMutex{}

	var waitGroup sync.WaitGroup
	waitGroup.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			for ref := range refChan {
				t.Logf("%v ref: %v", i, ref.String())
				count := expandRef(ref, azureRepoDir, t)

				lock.Lock()
				if count != 0 {
					res.Paths[ref.String()] = count
					res.AllProp += count
					res.AvailRef += 1
				}
				lock.Unlock()
			}
			waitGroup.Done()
		}(i)
	}

	for _, ref := range refList {
		refChan <- ref
	}
	close(refChan)

	waitGroup.Wait()

	t.Logf("total refs count: %d, ref with model: %d, total prop count: %d", res.AllRef, res.AvailRef, res.AllProp)

	b, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if testResultPath == "" {
		t.Log(string(b))
	} else {
		if err := os.WriteFile(testResultPath, b, 0644); err != nil {
			t.Fatal(err)
		}
	}

	return res
}

func getRefList(t *testing.T) []jsonreference.Ref {
	index, err := coverage.GetIndex()
	if err != nil {
		t.Fatal(err)
	}

	refMaps := make(map[string]jsonreference.Ref, 0)
	for resourceProvider, versionRaw := range index.ResourceProviders {
		for version, methodRaw := range versionRaw {
			for operationKind, resourceTypeRaw := range methodRaw {
				if operationKind != "PUT" && operationKind != "POST" {
					continue
				}
				for resourceType, operationInfo := range resourceTypeRaw {
					if operationInfo == nil {
						continue
					}
					for action, operationRefs := range operationInfo.Actions {
						for pathPattern, ref := range operationRefs {
							t.Logf("%s %s %s %s %s %s: %s", resourceProvider, version, operationKind, resourceType, action, pathPattern, ref.String())
							refMaps[ref.String()] = ref
						}
					}
					for pathPattern, ref := range operationInfo.OperationRefs {
						t.Logf("%s %s %s %s %s: %s", resourceProvider, version, operationKind, resourceType, pathPattern, ref.String())
						refMaps[ref.String()] = ref
					}
				}
			}
		}
	}

	refList := make([]jsonreference.Ref, 0)
	for _, ref := range refMaps {
		refList = append(refList, ref)
	}
	sort.Slice(refList, func(i, j int) bool {
		return refList[i].String() < refList[j].String()
	})

	for i, ref := range refList {
		t.Logf("refList ref %d: %s", i, ref.String())
	}

	return refList
}

func expandRef(ref jsonreference.Ref, azureRepoDir string, t *testing.T) int {
	model, err := coverage.GetModelInfoFromIndexRef(openapispec.Ref{Ref: ref}, azureRepoDir)
	if err != nil {
		panic(fmt.Errorf("get model info from index ref %s: %+v", ref.String(), err))
	}

	if model.ModelName == "" {
		t.Logf("model not found, skip %s", ref.String())
		return 0
	}

	expanded, err := coverage.Expand(model.ModelName, model.SwaggerPath)
	if err != nil {
		if strings.Contains(err.Error(), "not found in the definition of") {
			// https://github.com/Azure/azure-rest-api-specs/blob/f5cb37608399dd19760b9ef985a707294e32fbda/specification/vmware/resource-manager/Microsoft.AVS/stable/2021-06-01/vmware.json#L247
			t.Logf("model %s not found in the definition, skip %s", model.ModelName, ref.String())
			return 0
		}
		panic(fmt.Errorf("process %s, expand %s from %s: %+v", ref.String(), model.ModelName, model.SwaggerPath, err))
	}

	expanded.CountCoverage()
	return expanded.TotalCount
}
