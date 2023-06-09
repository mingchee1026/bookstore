package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Levap123/book_service/internal/book"
	"github.com/Levap123/book_service/internal/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepo struct {
	coll *mongo.Collection
	log  *logrus.Logger
}

func NewBookRepo(DB *mongo.Client, log *logrus.Logger) *BookRepo {
	return &BookRepo{
		coll: DB.Database("bookstore").Collection("books"),
		log:  log,
	}
}

func (br *BookRepo) Create(ctx context.Context, book book.Book) (string, error) {
	res, err := br.coll.InsertOne(ctx, book)
	if err != nil {
		return "", fmt.Errorf("book repo - create - %w", err)
	}

	ID := res.InsertedID.(primitive.ObjectID).Hex()
	return ID, nil
}

func (br *BookRepo) GetByID(ctx context.Context, bookID string) (book.Book, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return book.Book{}, fmt.Errorf("book repo - get object ID from hex - %w", err)
	}

	filter := bson.D{
		{"_id", objectID},
	}

	var bookBody book.Book
	if err := br.coll.FindOne(ctx, filter).Decode(&bookBody); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return book.Book{}, domain.ErrBookNotFound
		}
		return book.Book{}, fmt.Errorf("book repo - get one - %w", err)
	}

	return bookBody, err
}

func (br *BookRepo) GetAll(ctx context.Context) ([]book.Book, error) {
	cur, err := br.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("book repo - get all - %w", err)
	}

	books := make([]book.Book, 0, 10)

	for cur.Next(ctx) {
		var buffer book.Book
		if err := cur.Decode(&buffer); err != nil {
			return nil, fmt.Errorf("book repo - get all - %w", err)
		}

		books = append(books, buffer)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("book repo - get all - %w", err)
	}

	if len(books) == 0 {
		return nil, domain.ErrBookNotFound
	}

	return books, nil
}


func (br *BookRepo) Delete(ctx context.Context, bookID string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return "", fmt.Errorf("book repo - get object ID from hex - %w", domain.ErrBookNotFound)
	}

	filter := bson.D{
		{"_id", objectID},
	}

	res, err := br.coll.DeleteOne(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("book repo - delete by id - %w", err)
	}
	if res.DeletedCount != 1 {
		return "", domain.ErrBookNotFound
	}

	return bookID, nil
}

func (br *BookRepo) BooksFilter(ctx context.Context, genre, author, language, publisher []string) ([]book.Book, error) {
	filter := bson.M{}

	if len(genre) != 0 {
		filter["genre"] = bson.M{"$in": genre}
	}

	if len(author) != 0 {
		filter["author"] = bson.M{"$in": author}
	}

	if len(language) != 0 {
		filter["language"] = bson.M{"$in": language}
	}

	if len(publisher) != 0 {
		filter["publisher"] = bson.M{"$in": publisher}
	}

	return br.getByFilter(ctx, filter)
}

func (br *BookRepo) getByFilter(ctx context.Context, filter bson.M) ([]book.Book, error) {
	cur, err := br.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("book repo - get by filter - %w", err)
	}

	books := make([]book.Book, 0, 10)

	for cur.Next(ctx) {
		var buffer book.Book
		if err := cur.Decode(&buffer); err != nil {
			return nil, fmt.Errorf("book repo - decode - %w", err)
		}
		books = append(books, buffer)
	}

	return books, nil
}
