package models

type Channel struct {
	Label       string `xmlrpc:"label"`
	Name        string `xmlrpc:"name"`
	ParentLabel string `xmlrpc:"parent_label"`
	EndOfLife   string `xmlrpc:"end_of_life"`
	Arch        string `xmlrpc:"arch"`
}

type ChannelListItem struct {
	ID           int    `xmlrpc:"id"`
	Label        string `xmlrpc:"label"`
	Name         string `xmlrpc:"name"`
	ProviderName string `xmlrpc:"provider_name"`
	Packages     int    `xmlrpc:"packages"`
	Systems      int    `xmlrpc:"systems"`
	ArchName     string `xmlrpc:"arch_name"`
}

type ChannelDetails struct {
	ID               int    `xmlrpc:"id"`
	Label            string `xmlrpc:"label"`
	Name             string `xmlrpc:"name"`
	Summary          string `xmlrpc:"summary"`
	Description      string `xmlrpc:"description"`
	Arch             string `xmlrpc:"arch"`
	ParentChannel    string `xmlrpc:"parent_channel_label"`
	EndOfLife        string `xmlrpc:"end_of_life"`
	GpgCheck         bool   `xmlrpc:"gpg_check"`
	GpgKeyURL        string `xmlrpc:"gpg_key_url"`
	GpgKeyID         string `xmlrpc:"gpg_key_id"`
	Packages         int    `xmlrpc:"packages"`
	Errata           int    `xmlrpc:"errata"`
	ClonedFrom       string `xmlrpc:"cloned_from"`
	IsCloned         bool   `xmlrpc:"is_cloned"`
}

type ChannelPackage struct {
	ID      int    `xmlrpc:"id"`
	Name    string `xmlrpc:"name"`
	Version string `xmlrpc:"version"`
	Release string `xmlrpc:"release"`
	Epoch   string `xmlrpc:"epoch"`
	Arch    string `xmlrpc:"arch"`
}

type RepoInfo struct {
	ID               int      `xmlrpc:"id"`
	Label            string   `xmlrpc:"label"`
	SourceURL        string   `xmlrpc:"sourceUrl"`
	Type             string   `xmlrpc:"type"`
	HasSignedMetadata bool    `xmlrpc:"hasSignedMetadata"`
}
