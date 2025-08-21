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
	LinuxX86_64,
	LinuxArm64,
	MacosX86_64,
	MacosArm64,
}

var (
	LinuxX86_64 = Platform{Name: "linux", Arch: "x86_64"}
	LinuxArm64  = Platform{Name: "linux", Arch: "arm64"}
	MacosX86_64 = Platform{Name: "macos", Arch: "x86_64"}
	MacosArm64  = Platform{Name: "macos", Arch: "arm64"}
)
