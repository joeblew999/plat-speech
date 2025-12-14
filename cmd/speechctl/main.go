// speechctl - Stub CLI for speech artifacts (STT/TTS).
//
// This mirrors the codeccheck UX: install/download prebuilt assets, list, and check.
// Implementation is placeholder; wiring will hook into installer logic and manifests
// produced by .github/workflows/build-speech.yml.
package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "stt":
		handleStt(args)
	case "tts":
		handleTts(args)
	case "check":
		handleCheck(args)
	case "list":
		handleList(args)
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		fmt.Println("speechctl version", version)
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func handleStt(args []string) {
	fs := flag.NewFlagSet("stt", flag.ExitOnError)
	install := fs.Bool("install", false, "Download and install STT assets")
	dest := fs.String("dest", "./speech", "Destination directory for assets")
	device := fs.String("device", "auto", "Device preference: cpu|gpu|auto")
	model := fs.String("model", "", "Model selection (per release manifest)")
	manifest := fs.String("manifest", "", "Override manifest URL")
	offline := fs.Bool("offline", false, "Fail fast if network is required")
	upgrade := fs.Bool("upgrade", false, "Force reinstall even if present")
	jsonOut := fs.Bool("json", false, "Output JSON status")
	fs.Parse(args)

	if *install {
		fmt.Printf("[stub] STT install -> dest=%s device=%s model=%s manifest=%s offline=%v upgrade=%v\n",
			*dest, *device, *model, *manifest, *offline, *upgrade)
		return
	}

	fmt.Printf("[stub] STT status (json=%v)\n", *jsonOut)
}

func handleTts(args []string) {
	fs := flag.NewFlagSet("tts", flag.ExitOnError)
	install := fs.Bool("install", false, "Download and install TTS assets")
	dest := fs.String("dest", "./speech", "Destination directory for assets")
	device := fs.String("device", "auto", "Device preference: cpu|gpu|auto")
	voice := fs.String("voice", "", "Voice/style selection (per release manifest)")
	manifest := fs.String("manifest", "", "Override manifest URL")
	offline := fs.Bool("offline", false, "Fail fast if network is required")
	upgrade := fs.Bool("upgrade", false, "Force reinstall even if present")
	jsonOut := fs.Bool("json", false, "Output JSON status")
	fs.Parse(args)

	if *install {
		fmt.Printf("[stub] TTS install -> dest=%s device=%s voice=%s manifest=%s offline=%v upgrade=%v\n",
			*dest, *device, *voice, *manifest, *offline, *upgrade)
		return
	}

	fmt.Printf("[stub] TTS status (json=%v)\n", *jsonOut)
}

func handleCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "Output JSON status")
	fs.Parse(args)

	fmt.Printf("[stub] Check STT/TTS assets (json=%v)\n", *jsonOut)
}

func handleList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Parse(args)

	fmt.Println("[stub] List installed speech assets")
}

func printHelp() {
	fmt.Println("speechctl - manage speech assets (STT/TTS)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  speechctl stt [--install] [--dest DIR] [--device cpu|gpu|auto] [--model NAME] [--manifest URL] [--offline] [--upgrade] [--json]")
	fmt.Println("  speechctl tts [--install] [--dest DIR] [--device cpu|gpu|auto] [--voice NAME] [--manifest URL] [--offline] [--upgrade] [--json]")
	fmt.Println("  speechctl check [--json]")
	fmt.Println("  speechctl list")
	fmt.Println("  speechctl version")
	fmt.Println()
	fmt.Println("Commands mirror codeccheck: install prebuilt artifacts, list, and check.")
	fmt.Println("Assets are produced by the build-speech workflow and selected per OS/arch/flavor.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  speechctl stt --install --device auto")
	fmt.Println("  speechctl tts --install --voice F1")
	fmt.Println("  speechctl check --json")
	fmt.Println()
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Environment overrides:")
	fmt.Println("  SPEECH_MODELS_DIR, SPEECH_BIN_DIR, SPEECH_MANIFEST_URL (not yet implemented)")
	fmt.Println("  ONNXRUNTIME_LIB_PATH used by runners to locate bundled ORT")
	fmt.Println()
	fmt.Println("Note: implementation is pending; this stub preserves CLI shape for wiring installer logic.")
}
