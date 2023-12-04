package services

import (
	"context"
	"log"
	"subscriber/db"
	"subscriber/model"
	"subscriber/services/interfaces"
)

type PostgresOrderService struct {
	DB *db.PostgresDb
}

func NewPostgresOrderService(db *db.PostgresDb) interfaces.OrderService {
	return &PostgresOrderService{DB: db}
}

// Insert inserts one order to the db
func (orderService *PostgresOrderService) Insert(order *model.Order, ctx context.Context) error {
	conn, err := orderService.DB.Pool.Acquire(ctx)
	if err != nil {
		log.Println("cannot acquire a database connection", err)
		return err
	}
	defer conn.Release()

	// Insert into ORDERS
	query := `INSERT INTO ORDERS 
    (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = conn.Exec(ctx, query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OOFShard)
	if err != nil {
		log.Println("cannot insert into ORDERS", err)
		return err
	}

	// Insert into DELIVERY
	query = `INSERT into DELIVERY
	(order_id, name, phone, zip, city, address, region, email)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = conn.Exec(ctx, query, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Println("cannot insert into DELIVERY", err)
		return err
	}

	// Insert into PAYMENT
	query = `INSERT INTO PAYMENT
    (order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = conn.Exec(ctx, query, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Println("cannot insert into PAYMENT", err)
		return err
	}

	// Insert into ITEM
	for _, item := range order.Items {
		query = `insert into item
		(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = conn.Exec(ctx, query, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Println("cannot insert into ITEM", err)
			return err
		}
	}
	return nil
}

// InsertMany inserts many orders to the db
func (orderService *PostgresOrderService) InsertMany(order *model.Order, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// Get returns the order by id
func (orderService *PostgresOrderService) Get(id string, ctx context.Context, cache interfaces.Cache) (model.Order, error) {
	// trying to get order from the cache
	if value, isFound := cache.Get(id); isFound {
		order := value.(model.Order)
		log.Printf("%s is got from the cache\n", order.OrderUID)
		return order, nil
	}

	conn, err := orderService.DB.Pool.Acquire(ctx)
	if err != nil {
		log.Println("cannot acquire a database connection", err)
		return model.Order{}, err
	}
	defer conn.Release()

	var order model.Order

	// Get Order
	query := `SELECT 
    order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM ORDERS where order_uid=$1`
	err = conn.QueryRow(ctx, query, id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated, &order.OOFShard)
	if err != nil {
		log.Printf("cannot get from the db order with id=%s\n", id)
		return model.Order{}, err
	}

	// Get Delivery
	query = `SELECT 
    order_id, name, phone, zip, city, address, region, email
	FROM DELIVERY where order_id=$1`
	err = conn.QueryRow(ctx, query, id).Scan(&order.OrderUID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		log.Printf("cannot get from the db delivery with id=%s\n", id)
		return model.Order{}, err
	}

	// Get Payment
	query = `SELECT 
    order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM PAYMENT where order_id=$1`
	err = conn.QueryRow(ctx, query, id).Scan(&order.OrderUID, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		log.Printf("cannot get from the db payment with id=%s\n", id)
		return model.Order{}, err
	}

	// Get Items
	query = `SELECT 
    order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM ITEM where order_id=$1`
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Order{}, err
	}
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&order.OrderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Printf("cannot get from the db item with id=%s\n", id)
			return model.Order{}, err
		}
		order.Items = append(order.Items, item)
	}
	// push to the cache
	cache.Set(order.OrderUID, order, cache.GetDefaultExpiration())

	log.Printf("%s is got from the postgres\n", order.OrderUID)
	return order, nil
}

// GetMany returns {count} first entries from the db
func (orderService *PostgresOrderService) GetMany(count int, ctx context.Context) ([]model.Order, error) {
	var totalRowsCount int

	conn, err := orderService.DB.Pool.Acquire(ctx)
	if err != nil {
		log.Println("cannot acquire a database connection", err)
		return nil, err
	}
	defer conn.Release()
	err = conn.QueryRow(ctx, "SELECT count(*) FROM ORDERS").Scan(&totalRowsCount)
	if err != nil {
		log.Println("cannot get orders count", err)
		return []model.Order{}, nil
	}
	if totalRowsCount == 0 {
		return []model.Order{}, nil
	}

	rowsCount := count
	if count > totalRowsCount {
		rowsCount = totalRowsCount
	}

	orders := make(map[string]*model.Order, rowsCount)
	// Get Orders
	query := `SELECT 
    order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM ORDERS LIMIT $1`
	rows, err := conn.Query(ctx, query, rowsCount)
	if err != nil {
		log.Println("cannot query ORDERS")
		return []model.Order{}, err
	}
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated, &order.OOFShard)
		if err != nil {
			log.Println("cannot scan to order", err)
			return []model.Order{}, err
		}
		orders[order.OrderUID] = &order
	}

	// Get Delivery
	query = `SELECT 
    order_id, name, phone, zip, city, address, region, email
	FROM DELIVERY`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query DELIVERY")
		return []model.Order{}, err
	}
	for rows.Next() {
		var delivery model.Delivery
		var id string
		err := rows.Scan(&id, &delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
		if err != nil {
			log.Println("cannot scan to delivery", err)
			return []model.Order{}, err
		}
		if _, isFind := orders[id]; isFind {
			orders[id].Delivery = delivery
		}
	}

	// Get Payment
	query = `SELECT 
    order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM PAYMENT`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query PAYMENT")
		return []model.Order{}, err
	}
	for rows.Next() {
		var payment model.Payment
		var id string
		err := rows.Scan(&id, &payment.Transaction, &payment.RequestID, &payment.Currency,
			&payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
			&payment.GoodsTotal, &payment.CustomFee)
		if err != nil {
			log.Println("cannot scan to payment", err)
			return []model.Order{}, err
		}
		if _, isFind := orders[id]; isFind {
			orders[id].Payment = payment
		}
	}

	// Get Items
	query = `SELECT 
    order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM ITEM`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query ITEM")
		return []model.Order{}, err
	}
	for rows.Next() {
		var item model.Item
		var id string
		err := rows.Scan(&id, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Println("cannot scan to item", err)
			return []model.Order{}, err
		}
		if _, isFind := orders[id]; isFind {
			orders[id].Items = append(orders[id].Items, item)
		}
	}
	// convert map to slice
	ordersArray := make([]model.Order, rowsCount)
	curIndex := 0
	for _, value := range orders {
		ordersArray[curIndex] = *value
		curIndex++
	}
	return ordersArray, nil
}

// GetAll returns all entries from the db
func (orderService *PostgresOrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	var rowsCount int

	conn, err := orderService.DB.Pool.Acquire(ctx)
	if err != nil {
		log.Println("cannot acquire a database connection", err)
		return nil, err
	}
	defer conn.Release()
	err = conn.QueryRow(ctx, "SELECT count(*) FROM ORDERS").Scan(&rowsCount)
	if err != nil {
		log.Println("cannot get orders count", err)
		return []model.Order{}, nil
	}
	if rowsCount == 0 {
		return []model.Order{}, nil
	}

	orders := make(map[string]*model.Order, rowsCount)
	// Get Orders
	query := `SELECT 
    order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM ORDERS`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query ORDERS")
		return []model.Order{}, err
	}
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated, &order.OOFShard)
		if err != nil {
			log.Println("cannot scan to order", err)
			return []model.Order{}, err
		}
		orders[order.OrderUID] = &order
	}

	// Get Delivery
	query = `SELECT 
    order_id, name, phone, zip, city, address, region, email
	FROM DELIVERY`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query DELIVERY")
		return []model.Order{}, err
	}
	for rows.Next() {
		var delivery model.Delivery
		var id string
		err := rows.Scan(&id, &delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
		if err != nil {
			log.Println("cannot scan to delivery", err)
			return []model.Order{}, err
		}
		orders[id].Delivery = delivery
	}

	// Get Payment
	query = `SELECT 
    order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM PAYMENT`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query PAYMENT")
		return []model.Order{}, err
	}
	for rows.Next() {
		var payment model.Payment
		var id string
		err := rows.Scan(&id, &payment.Transaction, &payment.RequestID, &payment.Currency,
			&payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
			&payment.GoodsTotal, &payment.CustomFee)
		if err != nil {
			log.Println("cannot scan to payment", err)
			return []model.Order{}, err
		}
		orders[id].Payment = payment
	}

	// Get Items
	query = `SELECT 
    order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM ITEM`
	rows, err = conn.Query(ctx, query)
	if err != nil {
		log.Println("cannot query ITEM")
		return []model.Order{}, err
	}
	for rows.Next() {
		var item model.Item
		var id string
		err := rows.Scan(&id, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Println("cannot scan to item", err)
			return []model.Order{}, err
		}
		orders[id].Items = append(orders[id].Items, item)
	}
	// convert map to slice
	ordersArray := make([]model.Order, rowsCount)
	curIndex := 0
	for _, value := range orders {
		ordersArray[curIndex] = *value
		curIndex++
	}
	return ordersArray, nil
}
