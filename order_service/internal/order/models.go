package order

import (
	"time"

	"github.com/Levap123/order_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateOrderDTO struct {
	BookID string
	UserID uint64
}

type Order struct {
	ID      uint64    `db:"id"`
	BookID  string    `db:"book_id"`
	UserID  uint64    `db:"user_id"`
	AddedAt time.Time `db:"added_at"`
	Status  string    `db:"status"`
}

func createOrderDTOToOrder(dto CreateOrderDTO) Order {
	return Order{
		BookID:  dto.BookID,
		UserID:  dto.UserID,
		Status:  "создан",
		AddedAt: time.Now(),
	}
}

func fromReqToCreateDTO(req *proto.CreateOrderRequest) CreateOrderDTO {
	return CreateOrderDTO{
		BookID: req.BookId,
		UserID: req.UserId,
	}
}

func fromOrderToResp(order Order) *proto.Order {
	return &proto.Order{
		Id:      order.ID,
		UserId:  order.UserID,
		BookId:  order.BookID,
		Status:  order.Status,
		AddedAt: timestamppb.New(order.AddedAt),
	}
}
