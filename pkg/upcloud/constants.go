package upcloud

// Template UUIDs for different operating systems
// These can be found in UpCloud API documentation or via API calls
const (
	// Ubuntu templates
	TemplateUbuntu2204 = "01000000-0000-4000-8000-000020070100" // Ubuntu Server 22.04 LTS
	TemplateUbuntu2004 = "01000000-0000-4000-8000-000020060100" // Ubuntu Server 20.04 LTS

	// Debian templates
	TemplateDebian12 = "01000000-0000-4000-8000-000020050100" // Debian 12
	TemplateDebian11 = "01000000-0000-4000-8000-000020040100" // Debian 11

	// Rocky Linux templates
	TemplateRocky9 = "01000000-0000-4000-8000-000030090100" // Rocky Linux 9

	// AlmaLinux templates
	TemplateAlma9 = "01000000-0000-4000-8000-000040090100" // AlmaLinux 9
)

// DevPod status constants
const (
	StatusRunning  = "Running"
	StatusStopped  = "Stopped"
	StatusBusy     = "Busy"
	StatusNotFound = "NotFound"
)

// Default values
const (
	DefaultSSHUser     = "root"
	DefaultStorageTier = "maxiops"
	DefaultTimeout     = 300 // seconds
)

// Plan mappings - Legacy mapping for backward compatibility
// New plans are loaded from configs/server-plans.yaml
var PlanMap = map[string]string{
	// Developer Plans (New - Sept 2024)
	"DEV-1xCPU-1GB-10GB": "DEV-1xCPU-1GB-10GB",
	"DEV-1xCPU-1GB":      "DEV-1xCPU-1GB",
	"DEV-1xCPU-2GB":      "DEV-1xCPU-2GB",
	"DEV-1xCPU-4GB":      "DEV-1xCPU-4GB",
	"DEV-2xCPU-4GB":      "DEV-2xCPU-4GB",
	"DEV-2xCPU-8GB":      "DEV-2xCPU-8GB",
	"DEV-2xCPU-16GB":     "DEV-2xCPU-16GB",

	// Cloud Native Plans (New - Dec 2024)
	"CN-1xCPU-0.5GB": "CN-1xCPU-0.5GB",
	"CN-1xCPU-1GB":   "CN-1xCPU-1GB",
	"CN-2xCPU-2GB":   "CN-2xCPU-2GB",
	"CN-2xCPU-4GB":   "CN-2xCPU-4GB",

	// General Purpose Plans (Legacy)
	"1xCPU-1GB":   "1xCPU-1GB",
	"1xCPU-2GB":   "1xCPU-2GB",
	"2xCPU-4GB":   "2xCPU-4GB",
	"4xCPU-8GB":   "4xCPU-8GB",
	"6xCPU-16GB":  "6xCPU-16GB",
	"8xCPU-32GB":  "8xCPU-32GB",
	"12xCPU-48GB": "12xCPU-48GB",
	"16xCPU-64GB": "16xCPU-64GB",
	"20xCPU-96GB": "20xCPU-96GB",
}

// OS Image name to template UUID mapping
var ImageMap = map[string]string{
	"Ubuntu Server 22.04 LTS (Jammy Jellyfish)": TemplateUbuntu2204,
	"Ubuntu Server 20.04 LTS (Focal Fossa)":     TemplateUbuntu2004,
	"Debian 12 (Bookworm)":                      TemplateDebian12,
	"Debian 11 (Bullseye)":                      TemplateDebian11,
	"Rocky Linux 9":                             TemplateRocky9,
	"AlmaLinux 9":                               TemplateAlma9,
}
