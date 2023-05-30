## Armstrong Test Coverage

<blockquote><details open><summary>/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}</summary><blockquote>

<details><summary><span style="color:red">properties(3/339)</span></summary><blockquote>

- <span >description</span>

<details><summary><span style="color:red">outputs(2/337)</span></summary><blockquote>

<details><summary><span style="color:red">onError(0/2)</span></summary><blockquote>

- <span style="color:red">value=ContinueJob</span>

- <span style="color:red">value=StopProcessingJob</span>

</blockquote></details>

<details><summary><span style="color:red">preset(2/330)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

<details><summary><span style="color:red">#Microsoft.Media.AudioAnalyzerPreset(0/6)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">audioLanguage</span>

- <span style="color:red">experimentalOptions</span>

<details><summary><span style="color:red">mode(0/2)</span></summary><blockquote>

- <span style="color:red">value=Basic</span>

- <span style="color:red">value=Standard</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.BuiltInStandardEncoderPreset(2/31)</span></summary><blockquote>

- <span >@odata.type</span>

<details><summary><span style="color:red">configurations(0/12)</span></summary><blockquote>

- <span style="color:red">keyFrameIntervalInSeconds</span>

- <span style="color:red">maxBitrateBps</span>

- <span style="color:red">maxHeight</span>

- <span style="color:red">maxLayers</span>

- <span style="color:red">minBitrateBps</span>

- <span style="color:red">minHeight</span>

<details><summary><span style="color:red">complexity(0/3)</span></summary><blockquote>

- <span style="color:red">value=Balanced</span>

- <span style="color:red">value=Quality</span>

- <span style="color:red">value=Speed</span>

</blockquote></details>

<details><summary><span style="color:red">interleaveOutput(0/2)</span></summary><blockquote>

- <span style="color:red">value=InterleavedOutput</span>

- <span style="color:red">value=NonInterleavedOutput</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">presetName(1/17)</span></summary><blockquote>

- <span >value=AdaptiveStreaming</span>

- <span style="color:red">value=AACGoodQualityAudio</span>

- <span style="color:red">value=ContentAwareEncoding</span>

- <span style="color:red">value=ContentAwareEncodingExperimental</span>

- <span style="color:red">value=CopyAllBitrateNonInterleaved</span>

- <span style="color:red">value=DDGoodQualityAudio</span>

- <span style="color:red">value=H264MultipleBitrate1080p</span>

- <span style="color:red">value=H264MultipleBitrate720p</span>

- <span style="color:red">value=H264MultipleBitrateSD</span>

- <span style="color:red">value=H264SingleBitrate1080p</span>

- <span style="color:red">value=H264SingleBitrate720p</span>

- <span style="color:red">value=H264SingleBitrateSD</span>

- <span style="color:red">value=H265AdaptiveStreaming</span>

- <span style="color:red">value=H265ContentAwareEncoding</span>

- <span style="color:red">value=H265SingleBitrate1080p</span>

- <span style="color:red">value=H265SingleBitrate4K</span>

- <span style="color:red">value=H265SingleBitrate720p</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.FaceDetectorPreset(0/13)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">experimentalOptions</span>

<details><summary><span style="color:red">blurType(0/5)</span></summary><blockquote>

- <span style="color:red">value=Black</span>

- <span style="color:red">value=Box</span>

- <span style="color:red">value=High</span>

- <span style="color:red">value=Low</span>

- <span style="color:red">value=Med</span>

</blockquote></details>

<details><summary><span style="color:red">mode(0/3)</span></summary><blockquote>

- <span style="color:red">value=Analyze</span>

- <span style="color:red">value=Combined</span>

- <span style="color:red">value=Redact</span>

</blockquote></details>

<details><summary><span style="color:red">resolution(0/2)</span></summary><blockquote>

- <span style="color:red">value=SourceResolution</span>

- <span style="color:red">value=StandardDefinition</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.StandardEncoderPreset(0/269)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">experimentalOptions</span>

<details><summary><span style="color:red">codecs(0/170)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">label</span>

<details><summary><span style="color:red">#Microsoft.Media.AacAudio(0/9)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">bitrate</span>

- <span style="color:red">channels</span>

- <span style="color:red">label</span>

- <span style="color:red">samplingRate</span>

<details><summary><span style="color:red">profile(0/3)</span></summary><blockquote>

- <span style="color:red">value=AacLc</span>

- <span style="color:red">value=HeAacV1</span>

- <span style="color:red">value=HeAacV2</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.Audio(0/6)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">bitrate</span>

- <span style="color:red">channels</span>

- <span style="color:red">label</span>

- <span style="color:red">samplingRate</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.CopyAudio(0/3)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">label</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.CopyVideo(0/3)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">label</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.DDAudio(0/6)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">bitrate</span>

- <span style="color:red">channels</span>

- <span style="color:red">label</span>

- <span style="color:red">samplingRate</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.H264Video(0/41)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

- <span style="color:red">sceneChangeDetection</span>

<details><summary><span style="color:red">complexity(0/3)</span></summary><blockquote>

- <span style="color:red">value=Balanced</span>

- <span style="color:red">value=Quality</span>

- <span style="color:red">value=Speed</span>

</blockquote></details>

<details><summary><span style="color:red">layers(0/23)</span></summary><blockquote>

- <span style="color:red">adaptiveBFrame</span>

- <span style="color:red">bFrames</span>

- <span style="color:red">bitrate</span>

- <span style="color:red">bufferWindow</span>

- <span style="color:red">crf</span>

- <span style="color:red">frameRate</span>

- <span style="color:red">height</span>

- <span style="color:red">label</span>

- <span style="color:red">level</span>

- <span style="color:red">maxBitrate</span>

- <span style="color:red">referenceFrames</span>

- <span style="color:red">slices</span>

- <span style="color:red">width</span>

<details><summary><span style="color:red">entropyMode(0/2)</span></summary><blockquote>

- <span style="color:red">value=Cabac</span>

- <span style="color:red">value=Cavlc</span>

</blockquote></details>

<details><summary><span style="color:red">profile(0/6)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Baseline</span>

- <span style="color:red">value=High422</span>

- <span style="color:red">value=High444</span>

- <span style="color:red">value=High</span>

- <span style="color:red">value=Main</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">rateControlMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=ABR</span>

- <span style="color:red">value=CBR</span>

- <span style="color:red">value=CRF</span>

</blockquote></details>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.H265Video(0/33)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

- <span style="color:red">sceneChangeDetection</span>

<details><summary><span style="color:red">complexity(0/3)</span></summary><blockquote>

- <span style="color:red">value=Balanced</span>

- <span style="color:red">value=Quality</span>

- <span style="color:red">value=Speed</span>

</blockquote></details>

<details><summary><span style="color:red">layers(0/18)</span></summary><blockquote>

- <span style="color:red">adaptiveBFrame</span>

- <span style="color:red">bFrames</span>

- <span style="color:red">bitrate</span>

- <span style="color:red">bufferWindow</span>

- <span style="color:red">crf</span>

- <span style="color:red">frameRate</span>

- <span style="color:red">height</span>

- <span style="color:red">label</span>

- <span style="color:red">level</span>

- <span style="color:red">maxBitrate</span>

- <span style="color:red">referenceFrames</span>

- <span style="color:red">slices</span>

- <span style="color:red">width</span>

<details><summary><span style="color:red">profile(0/3)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Main10</span>

- <span style="color:red">value=Main</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.Image(0/14)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

- <span style="color:red">range</span>

- <span style="color:red">start</span>

- <span style="color:red">step</span>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.JpgImage(0/21)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

- <span style="color:red">range</span>

- <span style="color:red">spriteColumn</span>

- <span style="color:red">start</span>

- <span style="color:red">step</span>

<details><summary><span style="color:red">layers(0/6)</span></summary><blockquote>

- <span style="color:red">height</span>

- <span style="color:red">label</span>

- <span style="color:red">quality</span>

- <span style="color:red">width</span>

</blockquote></details>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.PngImage(0/19)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

- <span style="color:red">range</span>

- <span style="color:red">start</span>

- <span style="color:red">step</span>

<details><summary><span style="color:red">layers(0/5)</span></summary><blockquote>

- <span style="color:red">height</span>

- <span style="color:red">label</span>

- <span style="color:red">width</span>

</blockquote></details>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.Video(0/11)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">keyFrameInterval</span>

- <span style="color:red">label</span>

<details><summary><span style="color:red">stretchMode(0/3)</span></summary><blockquote>

- <span style="color:red">value=AutoFit</span>

- <span style="color:red">value=AutoSize</span>

- <span style="color:red">value=None</span>

</blockquote></details>

<details><summary><span style="color:red">syncMode(0/4)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=Cfr</span>

- <span style="color:red">value=Passthrough</span>

- <span style="color:red">value=Vfr</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">filters(0/62)</span></summary><blockquote>

<details><summary><span style="color:red">crop(0/5)</span></summary><blockquote>

- <span style="color:red">height</span>

- <span style="color:red">left</span>

- <span style="color:red">top</span>

- <span style="color:red">width</span>

</blockquote></details>

<details><summary><span style="color:red">deinterlace(0/6)</span></summary><blockquote>

<details><summary><span style="color:red">mode(0/2)</span></summary><blockquote>

- <span style="color:red">value=AutoPixelAdaptive</span>

- <span style="color:red">value=Off</span>

</blockquote></details>

<details><summary><span style="color:red">parity(0/3)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=BottomFieldFirst</span>

- <span style="color:red">value=TopFieldFirst</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">fadeIn(0/4)</span></summary><blockquote>

- <span style="color:red">duration</span>

- <span style="color:red">fadeColor</span>

- <span style="color:red">start</span>

</blockquote></details>

<details><summary><span style="color:red">fadeOut(0/4)</span></summary><blockquote>

- <span style="color:red">duration</span>

- <span style="color:red">fadeColor</span>

- <span style="color:red">start</span>

</blockquote></details>

<details><summary><span style="color:red">overlays(0/36)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">audioGainLevel</span>

- <span style="color:red">end</span>

- <span style="color:red">fadeInDuration</span>

- <span style="color:red">fadeOutDuration</span>

- <span style="color:red">inputLabel</span>

- <span style="color:red">start</span>

<details><summary><span style="color:red">#Microsoft.Media.AudioOverlay(0/8)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">audioGainLevel</span>

- <span style="color:red">end</span>

- <span style="color:red">fadeInDuration</span>

- <span style="color:red">fadeOutDuration</span>

- <span style="color:red">inputLabel</span>

- <span style="color:red">start</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.VideoOverlay(0/19)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">audioGainLevel</span>

- <span style="color:red">end</span>

- <span style="color:red">fadeInDuration</span>

- <span style="color:red">fadeOutDuration</span>

- <span style="color:red">inputLabel</span>

- <span style="color:red">opacity</span>

- <span style="color:red">start</span>

<details><summary><span style="color:red">cropRectangle(0/5)</span></summary><blockquote>

- <span style="color:red">height</span>

- <span style="color:red">left</span>

- <span style="color:red">top</span>

- <span style="color:red">width</span>

</blockquote></details>

<details><summary><span style="color:red">position(0/5)</span></summary><blockquote>

- <span style="color:red">height</span>

- <span style="color:red">left</span>

- <span style="color:red">top</span>

- <span style="color:red">width</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">rotation(0/6)</span></summary><blockquote>

- <span style="color:red">value=Auto</span>

- <span style="color:red">value=None</span>

- <span style="color:red">value=Rotate0</span>

- <span style="color:red">value=Rotate180</span>

- <span style="color:red">value=Rotate270</span>

- <span style="color:red">value=Rotate90</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">formats(0/34)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

<details><summary><span style="color:red">#Microsoft.Media.ImageFormat(0/3)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.JpgFormat(0/3)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.Mp4Format(0/7)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

<details><summary><span style="color:red">outputFiles(0/4)</span></summary><blockquote>

- <span style="color:red">labels(0/2)</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.MultiBitrateFormat(0/7)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

<details><summary><span style="color:red">outputFiles(0/4)</span></summary><blockquote>

- <span style="color:red">labels(0/2)</span>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.PngFormat(0/3)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.TransportStreamFormat(0/7)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">filenamePattern</span>

<details><summary><span style="color:red">outputFiles(0/4)</span></summary><blockquote>

- <span style="color:red">labels(0/2)</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">#Microsoft.Media.VideoAnalyzerPreset(0/9)</span></summary><blockquote>

- <span style="color:red">@odata.type</span>

- <span style="color:red">audioLanguage</span>

- <span style="color:red">experimentalOptions</span>

<details><summary><span style="color:red">insightsToExtract(0/3)</span></summary><blockquote>

- <span style="color:red">value=AllInsights</span>

- <span style="color:red">value=AudioInsightsOnly</span>

- <span style="color:red">value=VideoInsightsOnly</span>

</blockquote></details>

<details><summary><span style="color:red">mode(0/2)</span></summary><blockquote>

- <span style="color:red">value=Basic</span>

- <span style="color:red">value=Standard</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

<details><summary><span style="color:red">relativePriority(0/3)</span></summary><blockquote>

- <span style="color:red">value=High</span>

- <span style="color:red">value=Low</span>

- <span style="color:red">value=Normal</span>

</blockquote></details>

</blockquote></details>

</blockquote></details>

</blockquote></details>
</blockquote>
