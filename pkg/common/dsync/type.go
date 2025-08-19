package dsync

import (
	"context"
)

type Producer interface {
	Produce(c context.Context, args any, cn chan<- any) error
}
type Consumer interface {
	Consume(c context.Context, args any, ch <-chan any) error
}

type DataSaver func(args any, data any) error
type DataLoader func(args any) (data any, over bool, err error)

type SyncService interface {
	Sync(producerArgs, consumerArgs any) error
}

//type ReloadService interface {
//	ReloadLoader(args any) (data any, over bool, err error)
//	ReloadSaver(args any, data any) error
//}

//type ReloadInjector interface {
//	ReloadBefore(args *SyncArgs) error
//	ReloadAfter(args *SyncArgs) error
//}
