package domain

//nolint:revive // constant
var MAX_TOTAL_FILE_SIZE = 3 * 1024 * 1024 * 1024

type UserLimits struct {
	TotalFileSize int `bson:"total_file_size" json:"total_file_size"`
}
