package document

// Document is converts files of different formats (e.g CSV) to array of records (map[string]any)
// It assumes the first line of all files it process is the header
// It also writes the records keys as headers in the firstline of any file it outputs
