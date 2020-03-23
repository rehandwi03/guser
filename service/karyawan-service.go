package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rehandwi03/guser/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KaryawanService struct {
	db *sql.DB
}

func NewKaryawanService(db *sql.DB) model.KaryawanServiceServer {
	return &KaryawanService{db: db}
}

func (k *KaryawanService) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := k.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (k *KaryawanService) GetKaryawans(ctx context.Context, void *empty.Empty) (*model.AllKaryawans, error) {
	c, err := k.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	karyawans := []*model.Karyawan{}
	rows, err := c.QueryContext(ctx, "SELECT * FROM karyawan")
	if err != nil {
		log.Fatalf("could not exec the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		karyawan := new(model.Karyawan)
		if err := rows.Scan(&karyawan.Id, &karyawan.NamaLengkap, &karyawan.Alamat); err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		karyawans = append(karyawans, karyawan)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to retrieve data from user"+err.Error())
	}
	return &model.AllKaryawans{Karyawan: karyawans}, nil
}

func (k *KaryawanService) CreateKaryawan(ctx context.Context, req *model.CreateKaryawanRequest) (*model.CreateKaryawanResponse, error) {
	c, err := k.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "INSERT INTO karyawan VALUES(?,?)", req.Karyawan.NamaLengkap, req.Karyawan.Alamat)
	if err != nil {
		log.Fatalf("Cant insert data to table: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Id not found: %v", err)
	}
	karyawan := *&model.Karyawan{
		Id:          id,
		NamaLengkap: req.Karyawan.NamaLengkap,
		Alamat:      req.Karyawan.Alamat,
	}
	return &model.CreateKaryawanResponse{
		Karyawan: &karyawan,
	}, nil
}

func (k *KaryawanService) ReadKaryawan(ctx context.Context, req *model.GetKaryawanRequest) (*model.GetKaryawanResponse, error) {
	c, err := k.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT * FROM karyawan WHERE id = ?", req.Id)
	if err != nil {
		log.Fatalf("Cant exec the query: %v", err)
	}
	defer rows.Close()

	for !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "Cant fetch data"+err.Error())
		}
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Cant fetch data with id = '%d'", req.Id))
	}

	var karyawan model.Karyawan

	if err := rows.Scan(&karyawan.Id, &karyawan.NamaLengkap, &karyawan.Alamat); err != nil {
		return nil, status.Error(codes.Unknown, "Cant fetch data")
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Found multiple data with id = '%d'", req.Id))
	}
	return &model.GetKaryawanResponse{Karyawan: &karyawan}, nil
}

func (k *KaryawanService) UpdateKaryawan(ctx context.Context, req *model.UpdateKaryawanRequest) (*model.UpdateKaryawanResponse, error) {
	c, err := k.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE karyawan SET `nama_lengkap`, `alamat` WHERE `id` = ?", req.Karyawan.NamaLengkap, req.Karyawan.Alamat, req.Id)
	if err != nil {
		log.Fatalf("Can't exec the query: %v", err)
	}
	karyawan := &model.Karyawan{
		Id:          req.Id,
		NamaLengkap: req.Karyawan.NamaLengkap,
		Alamat:      req.Karyawan.Alamat,
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Karyawan with ID='%d' is not found",
			req.Id))
	}
	return &model.UpdateKaryawanResponse{Karyawan: karyawan}, nil
}

func (k *KaryawanService) DeleteKaryawan(ctx context.Context, req *model.DeleteKaryawanRequest) (*model.DeleteKaryawanResponse, error) {
	c, err := k.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM karyawan WHERE id = ?", req.Id)
	if err != nil {
		log.Fatalf("Cant exec the query: %v", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed deleting user"+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Data with id '%d' not found", req.Id))
	}
	return &model.DeleteKaryawanResponse{Message: "Deleted", Status: "204 No Content"}, nil
}
