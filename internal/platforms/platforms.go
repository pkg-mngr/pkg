package platforms

type Platform struct {
	Name string
	Arch string
}

func (p Platform) String() string {
	return p.Name + "-" + p.Arch
}

func ToPlatform(platform string) Platform {
	for _, p := range Platforms {
		if p.String() == platform {
			return p
		}
	}
	return Platform{}
}

var Platforms = []Platform{
	LinuxX64,
	LinuxArm64,
	MacosX64,
	MacosArm64,
}

var (
	LinuxX64   = Platform{Name: "linux", Arch: "x64"}
	LinuxArm64 = Platform{Name: "linux", Arch: "arm64"}
	MacosX64   = Platform{Name: "macos", Arch: "x64"}
	MacosArm64 = Platform{Name: "macos", Arch: "arm64"}
)
