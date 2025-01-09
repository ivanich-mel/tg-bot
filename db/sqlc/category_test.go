package db

import (
	"context"
	"database/sql"
	"melnik/telegram-bot/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}
func TestGetCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	category2, err := testQueries.GetCategory(context.Background(), category1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, category2)
	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.Name, category2.Name)
	require.Equal(t, category1.Balance, category2.Balance)
	require.WithinDuration(t, category1.CreatedAt, category2.CreatedAt, time.Second)
}

func TestUpdateCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	arg := UpdateCategoryParams{
		ID:      category1.ID,
		Name:    util.RandomName(),
		Balance: util.RandomBalance(),
	}
	err := testQueries.UpdateCategory(context.Background(), arg)

	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	require.Equal(t, arg.Name, category2.Name)
	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, arg.Balance, category2.Balance)
}

func TestDeleteCategory(t *testing.T) {
	category1 := createRandomCategory(t)
	err := testQueries.DeleteCategory(context.Background(), category1.ID)
	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
    require.Empty(t, category2)
}

func TestListCategories(t *testing.T) {
    for i := 0; i < 10; i++ {
        createRandomCategory(t)
    }

    arg := ListCategoriesParams{
        Limit:  5,
        Offset: 5,
    }
	categories, err := testQueries.ListCategories(context.Background(), arg)
    require.NoError(t, err)
    require.Len(t, categories, 5)

    for _, category := range categories {
        require.NotEmpty(t, category)
    }
}
func createRandomCategory(t *testing.T) Category {
	arg := CreateCategoryParams{
		Name:    util.RandomName(),
		Balance: util.RandomBalance(),
	}
	category, err := testQueries.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.Balance, category.Balance)
	require.NotZero(t, category.ID)
	require.NotZero(t, category.CreatedAt)
	return category
}
