package filesystem

import "github/pm/pkg/common"

type FileSystem struct {
	
}

func (fs *FileSystem) Update(alpha common.Alpha) (error) {
	alphaType := alpha.GetType();
	switch alphaType {
	case common.AddFileAlpha:
	case common.RemoveFileAlpha:
	}
	return nil
}

func (fs *FileSystem) Rewind(alpha common.Alpha) (error) {
	alphaType := alpha.GetType();
	switch alphaType {
	case common.AddFileAlpha:
	case common.RemoveFileAlpha:
	}
	return nil
}

func (fs *FileSystem) Validate(alpha common.Alpha) (error) {
	alphaType := alpha.GetType();
	switch alphaType {
	case common.AddFileAlpha:
	case common.RemoveFileAlpha:
	}
	return nil
}

type AddFileAlpha struct {
}

func (afa AddFileAlpha) GetType() (byte) {
	return common.AddFileAlpha
}

func (afa AddFileAlpha) GetId() (string) {
	return ""
}

func (afa AddFileAlpha) GetHash() (string) {
	return ""
}

func (afa AddFileAlpha) SetHash(alpha common.Alpha) {
}


type RemoveFileAlpha struct {
}

func (rfa RemoveFileAlpha) GetType() (byte) {
	return common.RemoveFileAlpha
}

func (rfa RemoveFileAlpha) GetId() (string) {
	return ""
}

func (rfa RemoveFileAlpha) GetHash() (string) {
	return ""
}

func (rfa RemoveFileAlpha) SetHash() {
}


