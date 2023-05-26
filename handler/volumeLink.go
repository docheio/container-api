package handler

type VolumeLink struct {
	Name  string
	Mount string
	Claim string
}

type VolumeLinkGroup []VolumeLink

type VolumeLinks []VolumeLinkGroup
