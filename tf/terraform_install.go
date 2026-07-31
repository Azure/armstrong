package tf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
)

const terraformBinary = "terraform"

// FindTerraform finds the path to the terraform executable whose version meets the min/max version constraint.
// It first tries to find from the local OS PATH. If there is no match, it will then download the release of the minVersion from hashicorp to the tfDir.
func FindTerraform(ctx context.Context) (string, error) {

	// Initialize the workspace
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("error finding the user cache directory: %w", err)
	}
	rootDir := filepath.Join(cacheDir, "armstrong")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return "", fmt.Errorf("creating workspace root %q: %w", rootDir, err)
	}
	tfDir := filepath.Join(rootDir, "terraform")
	if err := os.MkdirAll(tfDir, 0755); err != nil {
		return "", fmt.Errorf("creating terraform cache dir %q: %w", tfDir, err)
	}

	installer := install.NewInstaller()
	terraformPath, err := installer.Ensure(ctx, []src.Source{
		&fs.AnyVersion{ExactBinPath: filepath.Join(tfDir, terraformBinary)},
		&fs.AnyVersion{Product: &product.Terraform},
		&releases.LatestVersion{Product: product.Terraform, InstallDir: tfDir},
	})
	if err != nil {
		return "", fmt.Errorf("finding or installing terraform: %w", err)
	}

	return terraformPath, nil
}
