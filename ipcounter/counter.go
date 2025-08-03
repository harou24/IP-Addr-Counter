package ipcounter

// Counter defines a method to count unique IPs from a file.
type Counter interface {
	CountUniqueIPs(filename string) (int64, error)
}
