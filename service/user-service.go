package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rehandwi03/guser/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	db *sql.DB
}

func NewUserServiceServer(db *sql.DB) model.UserServiceServer {
	return &UserService{db: db}
}

func (u *UserService) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := u.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (u *UserService) GetUsers(ctx context.Context, void *empty.Empty) (*model.AllUsers, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get semua data
	users := []*model.User{}
	rows, err := c.QueryContext(ctx, "SELECT id, username, password from user")
	if err != nil {
		log.Fatalf("could not get data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := new(model.User)

		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			return nil, status.Errorf(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to retrieve data from user"+err.Error())
	}
	return &model.AllUsers{Message: "Ok", Status: "200 OK", User: users}, nil
}

func (u *UserService) Create(ctx context.Context, req *model.CreateRequest) (*model.CreateResponse, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// log.Printf("Data: %v", req.Item.Username)
	res, err := c.ExecContext(ctx, "INSERT INTO user (`username`,`password`,`karyawan_id`) VALUES(?, ?, ?)", req.Item.Username, req.Item.Password, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert data user"+err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Error())
	}
	user := *&model.User{
		Id:       id,
		Username: req.Item.Username,
	}
	return &model.CreateResponse{
		Status:  200,
		Message: "Successfully",
		User:    &user,
	}, nil
}

func (u *UserService) Read(ctx context.Context, req *model.ReadRequest) (*model.ReadResponse, error) {
	c, err := u.connect(ctx)
	if err != nil {
		log.Fatalf("Cant connect to db: %v", err)
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT * FROM user WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Cant fetch data")
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "Cant fetch data")
		}
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Cant find data with id = '%d'", req.Id))
	}

	var user model.User
	var karyawan_id = struct {
		karyawan_id int64
	}{}
	if err := rows.Scan(&user.Id, &user.Username, &user.Password, &karyawan_id.karyawan_id); err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Cant fetch data", err))
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Found multiple data with id = '%d'", req.Id))
	}

	// ambil data ke karyawan lewat grpc
	conn, err := grpc.Dial("karyawan-svc:7776", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	k := model.NewKaryawanServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request := model.GetKaryawanRequest{
		Id: karyawan_id.karyawan_id,
	}
	res, err := k.ReadKaryawan(ctx, &request)
	if err != nil {
		log.Fatalf("cant get data from karyawan service: %v", err)
	}
	log.Printf("Data from karyawan: %v", res)

	karyawan := new(model.Karyawan)
	// karyawan.Id = res.Karyawan.Id
	karyawan.NamaLengkap = res.Karyawan.NamaLengkap
	karyawan.Alamat = res.Karyawan.Alamat
	return &model.ReadResponse{
		User:     &user,
		Karyawan: karyawan,
	}, nil
}

func (u *UserService) Update(ctx context.Context, req *model.UpdateRequest) (*model.UpdateResponse, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE user SET `username` = ?, `password` = ? WHERE `id` = ?", req.User.Username, req.User.Password, req.User.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}
	user := &model.User{
		Id:       req.User.Id,
		Username: req.User.Username,
		Password: req.User.Password,
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("User with ID='%d' is not found",
			req.User.Id))
	}

	return &model.UpdateResponse{Message: "Updated", Status: "200 OK", User: user}, nil
}

func (u *UserService) Delete(ctx context.Context, req *model.DeleteRequest) (*model.DeleteResponse, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM user WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to delete user"+err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	return &model.DeleteResponse{
		Message: "Delete Success",
		Status:  "204",
	}, nil
}
