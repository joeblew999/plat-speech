# Local STT/TTS plan

## Objective
- Deliver fully local, install-once STT + TTS on macOS/Linux/Windows (plus Pi/e-reader where feasible) with a Go-first structure: `speechctl`-style CLI (patterned after `cmd/codeccheck`) and GitHub workflows that prebuild artifacts users pull—no manual Git LFS needed for end users.

- **Stacks:**  
  - **STT:** sherpa-onnx (ONNX Runtime backends; streaming; Go API examples in `.src/sherpa-onnx/go-api-examples`).  
  - **TTS:** Supertonic (ONNX/WebGPU/WASM; cloned in `.src/supertonic`).

## Code analysis targets (to mirror and wrap)
- **STT (sherpa-onnx):**
  - Streaming decode example: `.src/sherpa-onnx/go-api-examples/streaming-decode-files/main.go` (OnlineRecognizer, flag-driven model config, 16 kHz mono PCM handling).
  - Mic/real-time: `.src/sherpa-onnx/go-api-examples/real-time-speech-recognition-from-microphone`.
  - VAD + ASR combos: `vad-asr-paraformer`, `vad-asr-whisper`.
  - Takeaways to wrap: config structs, provider selection (`cpu/cuda/coreml/directml`), streaming API lifecycle, WAV ingestion.
- **TTS (Supertonic):**
  - Go runner: `.src/supertonic/go/example_onnx.go` + `helper.go` (ONNX Runtime binding, voice presets under `go/assets/voice_styles/`).
  - Dependency: imports `github.com/yalue/onnxruntime_go v1.11.0` (see `.src/supertonic/go/go.mod`), so ORT versioning must match that wrapper; it’s CGO-based and expects a matching ORT shared library.
  - Notes: needs ONNX Runtime shared lib unless we bundle one; batch/long-form options; voice/style selection via JSON presets.

## Delivery & installer (mirror `codeccheck`)
- **CLI name:** `speechctl` (Go), placed under `cmd/speechctl/`.
- **Commands:** `stt install`, `tts install`, `check`, `list`, `--json`.
- Auto-detect OS/arch/GPU; download the correct speech artifact from releases; verify SHA256; expand into `bin/`, `models/stt/`, `models/tts/`.
- Flags: `--dest`, `--device cpu|gpu|auto`, `--model`, `--voice`, `--manifest-url`, `--offline`, `--upgrade`.
- No Git LFS for users: artifacts are prepacked zips/tars with models+binaries+manifest+README.

## ONNX Runtime handling
- Both sherpa-onnx and Supertonic require ONNX Runtime. Bundle per-platform ORT binaries in artifacts to avoid user installs.
- Variants: CPU-only everywhere; optional CUDA build for Linux/Windows; CoreML/Metal/DirectML where available. Pick one CPU build as default; include GPU builds only where tested.
- Include ORT version + hashes in `manifest.json`; set loader/env hints in `speechctl` (e.g., `ONNXRUNTIME_LIB_PATH`) when running smoke tests.
- Upstream watch: yzma#130 proposes purego/ffi bindings for ONNX Runtime (no CGO). If it lands, we can drop bundled shared libs and ship a single Go binary, improving cross-compile and scratch-container support.
- Related work: yalue/onnxruntime_go already maps the full ORT C API (tensors, sessions, CUDA/DirectML providers) but uses CGO. If purego is not upstreamed soon, we could adapt that surface, swapping CGO calls for purego/ffi to cut effort.

## Platform/ISA matrix (artifacts + installer logic)
- **macOS arm64:** bundle ORT CPU; optional Metal/CoreML build if verified with sherpa/Supertonic. `speechctl` prefers CPU unless GPU artifact present and `--device gpu|auto`.
- **Linux amd64:** bundle ORT CPU (glibc); optional CUDA build (match tested CUDA version). `speechctl` defaults CPU, uses CUDA if artifact + driver present.
- **Linux arm64 (server/Jetson):** bundle ORT CPU; optional CUDA if tested for the target. Same selection logic.
- **Raspberry Pi arm64:** CPU-only ORT; smallest sherpa model + light Supertonic voices; no GPU path.
- **Windows amd64:** ORT CPU; optional DirectML if tested. Default to CPU; allow `--device gpu` to pick DML artifact.
- **E-reader class:** treat like Pi (CPU-only, smallest assets).
- Installer selection: detect OS/arch, look for matching GPU artifact; fall back to CPU artifact if none or if probe fails. Record chosen artifact in manifest.
- Versioning: lock ORT version per release (aligned with yalue/onnxruntime_go v1.11.0 unless we vendor purego); include SHA256 per platform in manifest.
## CI workflow outline (like `build-codecs.yml`)
- **Workflow file:** `.github/workflows/build-speech.yml`.
- Matrix: macOS, Windows, Linux; amd64/arm64; add Pi arm64 if runner available.
- Steps per job: checkout; install git-lfs; fetch Supertonic models/voices; fetch sherpa-onnx models (streaming conformer/paraformer, light variants for Pi/e-reader); fetch/build ORT (CPU + optional GPU variant per matrix); cache model/ORT downloads; run smoke tests (short STT clip, short TTS phrase); package `speech-{os}-{arch}.zip`; upload artifacts.
- Artifacts layout: `bin/` (runners/wrappers), `models/stt/`, `models/tts/`, `manifest.json` (hashes/sizes), `README` (usage), license snippets if required.
- Caching: use `actions/cache` keyed by ORT version + platform to speed rebuilds; cache model downloads where HF terms allow; avoid redundant git-lfs pulls across jobs if possible.
- Naming: release artifacts as `speech-${os}-${arch}-${flavor}.zip` (e.g., `speech-linux-amd64-cpu.zip`, `speech-linux-amd64-cuda.zip`, `speech-macos-arm64-cpu.zip`, `speech-windows-amd64-dml.zip`, `speech-linux-arm64-pi-cpu.zip`) to map cleanly in `speechctl` selection logic.

## Integration notes (mediaDevices desktop stack)
- **STT capture:** tap mic before network send; resample to 16 kHz mono; light VAD; small chunks with overlap; target <1–2s latency for captions. Feed sherpa-onnx streaming API; store timestamps for UI alignment.
- **TTS playback:** synthesize via Supertonic; inject as a separate outbound audio track so users can mute it; match sample rate to the existing audio graph.
- **Config:** per-endpoint device IDs; device preference `cpu|gpu|auto`; model/voice paths; sensible desktop defaults; allow offline mode (no downloads during calls).

## Platform notes
- **Raspberry Pi:** use light sherpa-onnx streaming models; small chunk sizes; aim for RTF ≤1.5×. Supertonic with light presets; keep total model footprint to a few hundred MB; watch thermals.
- **E-reader:** TTS via Supertonic with lowest-footprint voices at 16–22.05 kHz; STT limited to short commands with very small models; e-ink-friendly UI; bundle models, no runtime downloads.

## References
- STT examples: `.src/sherpa-onnx/go-api-examples` (use as pattern for smoke tests and API surface).
- TTS examples: `.src/supertonic` (ONNX/WebGPU/WASM; pick minimal runners for packaging).
