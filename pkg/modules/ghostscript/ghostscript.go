package ghostscript

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/gotenberg/gotenberg/v8/pkg/gotenberg"
)

func init() {
	gotenberg.MustRegisterModule(new(Ghostscript))
}

// Ghostscript abstracts the CLI tool Ghostscript and implements the [gotenberg.PdfEngine]
// interface.
type Ghostscript struct {
	binPath string
	libPath string
	iccRgbPath string
}

// Descriptor returns a [Ghostscript]'s module descriptor.
func (engine *Ghostscript) Descriptor() gotenberg.ModuleDescriptor {
	return gotenberg.ModuleDescriptor{
		ID:  "ghostscript",
		New: func() gotenberg.Module { return new(Ghostscript) },
	}
}

// Provision sets the module properties.
func (engine *Ghostscript) Provision(ctx *gotenberg.Context) error {
	binPath, ok := os.LookupEnv("GHOSTSCRIPT_BIN_PATH")
	if !ok {
		return errors.New("GHOSTSCRIPT_BIN_PATH environment variable is not set")
	}

	libPath, ok := os.LookupEnv("GHOSTSCRIPT_LIB_PATH")
	if !ok {
		return errors.New("GHOSTSCRIPT_LIB_PATH environment variable is not set")
	}

	iccRgbPath, ok := os.LookupEnv("GHOSTSCRIPT_ICC_RGB_PATH")
	if !ok {
		return errors.New("GHOSTSCRIPT_ICC_RGB_PATH environment variable is not set")
	}

	engine.binPath = binPath
	engine.libPath = libPath
	engine.iccRgbPath = iccRgbPath

	return nil
}

// Validate validates the module properties.
func (engine *Ghostscript) Validate() error {
	_, err := os.Stat(engine.binPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Ghostscript binary path does not exist: %w", err)
	}

	return nil
}

// Metrics returns the metrics.
// func (engine Ghostscript) Metrics() ([]gotenberg.Metric, error) {
// 	return []gotenberg.Metric{
// 		{
// 			Name:        "ghostscript_active_instances_count",
// 			Description: "Current number of active Ghostscript instances.",
// 			Read: func() float64 {
// 				activeInstancesCountMu.RLock()
// 				defer activeInstancesCountMu.RUnlock()
//
// 				return activeInstancesCount
// 			},
// 		},
// 	}, nil
// }

// Split splits a given PDF file.
func (engine *Ghostscript) Split(ctx context.Context, logger *slog.Logger, mode gotenberg.SplitMode, inputPath, outputDirPath string) ([]string, error) {
	return nil, fmt.Errorf("split PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Merge combines multiple PDFs into a single PDF.
func (engine *Ghostscript) Merge(ctx context.Context, logger *slog.Logger, inputPaths []string, outputPath string) error {
	var args []string
    args = append(args, "-dBATCH")
    args = append(args, "-dNOPAUSE")
    args = append(args, "-sDEVICE=pdfwrite")
    args = append(args, fmt.Sprintf("-sOutputFile=%s", outputPath))
    args = append(args, inputPaths...)

    cmd, err := gotenberg.CommandContext(ctx, logger, engine.binPath, args...)
    if err != nil {
        return fmt.Errorf("create command: %w", err)
    }

//     activeInstancesCountMu.Lock()
//     activeInstancesCount += 1
//     activeInstancesCountMu.Unlock()

    _, err = cmd.Exec()

//     activeInstancesCountMu.Lock()
//     activeInstancesCount -= 1
//     activeInstancesCountMu.Unlock()

    if err == nil {
        return nil
    }

    return fmt.Errorf("merge PDFs with Ghostscript: %w", err)
}

// Flatten is not available in this implementation.
func (engine *Ghostscript) Flatten(ctx context.Context, logger *slog.Logger, inputPath string) error {
	return fmt.Errorf("flatten PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Convert is not available in this implementation.
func (engine *Ghostscript) Convert(ctx context.Context, logger *slog.Logger, formats gotenberg.PdfFormats, inputPath, outputPath string) error {
	var pdfALevel string

    switch formats.PdfA {
    case gotenberg.PdfA1b:
        pdfALevel = "1"
    case gotenberg.PdfA2b:
        pdfALevel = "2"
    case gotenberg.PdfA3b:
        pdfALevel = "3"
    default:
        return fmt.Errorf("convert PDF to '%s' with ghostsript: %w", formats.PdfA, gotenberg.ErrPdfFormatNotSupported)
    }

    var args []string
    args = append(args, fmt.Sprintf("-dPDFA=%s", pdfALevel))
    args = append(args, "-dBATCH")
    args = append(args, "-dNOPAUSE")
    args = append(args, "-sDEVICE=pdfwrite")
    args = append(args, "-sColorConversionStrategy=RGB")
    args = append(args, "-dPDFACompatibilityPolicy=1")
    args = append(args, fmt.Sprintf("--permit-file-read=%s", engine.iccRgbPath))
    args = append(args, fmt.Sprintf("-sOutputFile=%s", outputPath))
    args = append(args, fmt.Sprintf("%s/PDFA_def.ps", engine.libPath))
    args = append(args, inputPath)

    cmd, err := gotenberg.CommandContext(ctx, logger, engine.binPath, args...)
    if err != nil {
        return fmt.Errorf("create command: %w", err)
    }

//     activeInstancesCountMu.Lock()
//     activeInstancesCount += 1
//     activeInstancesCountMu.Unlock()

    _, err = cmd.Exec()

//     activeInstancesCountMu.Lock()
//     activeInstancesCount -= 1
//     activeInstancesCountMu.Unlock()

    if err == nil {
        return nil
    }
    return fmt.Errorf("convert PDF to '%s' with Ghostscript: %w", formats.PdfA, err)
}

// ReadMetadata is not available in this implementation.
func (engine *Ghostscript) ReadMetadata(ctx context.Context, logger *slog.Logger, inputPath string) (map[string]any, error) {
	return nil, fmt.Errorf("read PDF metadata with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// WriteMetadata is not available in this implementation.
func (engine *Ghostscript) WriteMetadata(ctx context.Context, logger *slog.Logger, metadata map[string]any, inputPath string) error {
	return fmt.Errorf("write PDF metadata with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Encrypt adds password protection to a PDF file using Ghostscript.
func (engine *Ghostscript) Encrypt(ctx context.Context, logger *slog.Logger, inputPath string, opts gotenberg.EncryptOptions) error {
	return fmt.Errorf("encrypt PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// EmbedFiles is not available in this implementation.
func (engine *Ghostscript) EmbedFiles(ctx context.Context, logger *slog.Logger, filePaths []string, inputPath string) error {
	return fmt.Errorf("embed files with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// EmbedFilesMetadata is not available in this implementation.
func (engine *Ghostscript) EmbedFilesMetadata(ctx context.Context, logger *slog.Logger, metadata map[string]map[string]string, inputPath string) error {
	return fmt.Errorf("embed files metadata with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// OptimizeImages is not available in this implementation.
func (engine *Ghostscript) OptimizeImages(ctx context.Context, logger *slog.Logger, imageQuality int, inputPath string) error {
	return fmt.Errorf("optimize images with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// PageCount is not available in this implementation.
func (engine *Ghostscript) PageCount(ctx context.Context, logger *slog.Logger, inputPath string) (int, error) {
	return 0, fmt.Errorf("count PDF pages with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// ReadBookmarks is not available in this implementation.
func (engine *Ghostscript) ReadBookmarks(ctx context.Context, logger *slog.Logger, inputPath string) ([]gotenberg.Bookmark, error) {
	return nil, fmt.Errorf("read PDF bookmarks with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// WriteBookmarks is not available in this implementation.
func (engine *Ghostscript) WriteBookmarks(ctx context.Context, logger *slog.Logger, inputPath string, bookmarks []gotenberg.Bookmark) error {
	return fmt.Errorf("write PDF bookmarks with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Watermark is not available in this implementation.
func (engine *Ghostscript) Watermark(ctx context.Context, logger *slog.Logger, inputPath string, stamp gotenberg.Stamp) error {
	return fmt.Errorf("watermark PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Stamp is not available in this implementation.
func (engine *Ghostscript) Stamp(ctx context.Context, logger *slog.Logger, inputPath string, stamp gotenberg.Stamp) error {
	return fmt.Errorf("stamp PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// Rotate is not available in this implementation.
func (engine *Ghostscript) Rotate(ctx context.Context, logger *slog.Logger, inputPath string, angle int, pages string) error {
	return fmt.Errorf("rotate PDF with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// InjectFacturXXMP is not available in this implementation.
func (engine *Ghostscript) InjectFacturXXMP(ctx context.Context, logger *slog.Logger, facturX gotenberg.FacturX, inputPath string) error {
	return fmt.Errorf("inject Factur-X XMP with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// ReadPdfAConformance is not available in this implementation.
func (engine *Ghostscript) ReadPdfAConformance(ctx context.Context, logger *slog.Logger, inputPath string) (string, string, error) {
	return "", "", fmt.Errorf("read PDF/A conformance with Ghostscript: %w", gotenberg.ErrPdfEngineMethodNotSupported)
}

// var (
// 	activeInstancesCount   float64
// 	activeInstancesCountMu sync.RWMutex
// )

// Interface guards.
var (
	_ gotenberg.Module           = (*Ghostscript)(nil)
	_ gotenberg.Provisioner      = (*Ghostscript)(nil)
	_ gotenberg.Validator        = (*Ghostscript)(nil)
// 	_ gotenberg.MetricsProvider  = (*Ghostscript)(nil)
	_ gotenberg.PdfEngine        = (*Ghostscript)(nil)
)