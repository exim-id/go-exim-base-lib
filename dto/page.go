package dto

type PageReq[T interface{}] struct {
	Page  uint64
	Size  uint64
	Where T
}

type PageRes[T interface{}] struct {
	Page      uint64
	Size      uint64
	PageCount uint64
	Datas     []T
}
