package models

import "time"

type System struct {
	ID       int    `xmlrpc:"id"`
	Name     string `xmlrpc:"name"`
	Hostname string `xmlrpc:"hostname"`
}

type SystemDetails struct {
	ID                int      `xmlrpc:"id"`
	ProfileName       string   `xmlrpc:"profile_name"`
	MachineID         string   `xmlrpc:"machine_id"`
	MinionID          string   `xmlrpc:"minion_id"`
	BaseEntitlement   string   `xmlrpc:"base_entitlement"`
	AddonEntitlements []string `xmlrpc:"addon_entitlements"`
	AutoUpdate        bool     `xmlrpc:"auto_update"`
	Release           string   `xmlrpc:"release"`
	Address1          string   `xmlrpc:"address1"`
	Address2          string   `xmlrpc:"address2"`
	City              string   `xmlrpc:"city"`
	State             string   `xmlrpc:"state"`
	Country           string   `xmlrpc:"country"`
	Hostname          string   `xmlrpc:"hostname"`
	LastBoot          string   `xmlrpc:"last_boot"`
	OSAStatus         string   `xmlrpc:"osa_status"`
	LockStatus        bool     `xmlrpc:"lock_status"`
	Virtualization    string   `xmlrpc:"virtualization"`
	ContactMethod     string   `xmlrpc:"contact_method"`
	Description       string   `xmlrpc:"description"`
	Payg              bool     `xmlrpc:"payg"`
}

type CPUInfo struct {
	Cache       string `xmlrpc:"cache"`
	Family      string `xmlrpc:"family"`
	MHz         string `xmlrpc:"mhz"`
	Flags       string `xmlrpc:"flags"`
	Model       string `xmlrpc:"model"`
	Vendor      string `xmlrpc:"vendor"`
	Arch        string `xmlrpc:"arch"`
	Stepping    string `xmlrpc:"stepping"`
	Count       string `xmlrpc:"count"`
	SocketCount string `xmlrpc:"socket_count"`
	CoreCount   string `xmlrpc:"core_count"`
	ThreadCount string `xmlrpc:"thread_count"`
}

type MemoryInfo struct {
	RAM  string `xmlrpc:"ram"`
	Swap string `xmlrpc:"swap"`
}

type DMIInfo struct {
	Vendor      string `xmlrpc:"vendor"`
	System      string `xmlrpc:"system"`
	Product     string `xmlrpc:"product"`
	Asset       string `xmlrpc:"asset"`
	Board       string `xmlrpc:"board"`
	BiosRelease string `xmlrpc:"bios_release"`
	BiosVendor  string `xmlrpc:"bios_vendor"`
	BiosVersion string `xmlrpc:"bios_version"`
}

type InstalledPackage struct {
	Name    string `xmlrpc:"name"`
	Version string `xmlrpc:"version"`
	Release string `xmlrpc:"release"`
	Arch    string `xmlrpc:"arch"`
	Epoch   string `xmlrpc:"epoch"`
}

type EventHistory struct {
	ID         int    `xmlrpc:"id"`
	ActionName string `xmlrpc:"action_name"`
	Status     string `xmlrpc:"status"`
	Completed  string `xmlrpc:"completed"`
	Message    string `xmlrpc:"message"`
}

type SystemGroup struct {
	ID          int    `xmlrpc:"id"`
	Name        string `xmlrpc:"name"`
	Description string `xmlrpc:"description"`
}

type EventHistoryResult struct {
	EarliestDate *time.Time `xmlrpc:"earliestDate"`
	Offset       int        `xmlrpc:"offset"`
	Limit        int        `xmlrpc:"limit"`
	TotalSize    int        `xmlrpc:"total_size"`
	Events       []EventHistory
}
