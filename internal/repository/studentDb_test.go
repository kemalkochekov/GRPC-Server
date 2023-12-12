//go:build integration
// +build integration

package repository

import (
	"GRPC_Server/internal/repository/postgres"
	"GRPC_Server/internal/server/serviceEntities"
	"GRPC_Server/pkg/pkgErrors"
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestCreateStudent(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)
	})
}
func TestGetStudent(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)
		respStudentGet, err := studentRepo.GetByID(ctx, respStudent)
		//assert
		require.NoError(t, err)
		assert.Equal(t, testStudentReq.StudentName, respStudentGet.StudentName)
		assert.Equal(t, testStudentReq.Grade, respStudentGet.Grade)
	})
	t.Run("Fail", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudentGet, err := studentRepo.GetByID(ctx, testStudentReq.StudentID)
		//assert
		require.Error(t, err)
		assert.Equal(t, serviceEntities.StudentRequest{}, respStudentGet)
	})
}
func TestUpdateStudent(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)
		//arrange
		testStudentReq.StudentName = "Test2"
		testStudentReq.Grade = 92
		//act
		err = studentRepo.Update(ctx, respStudent, testStudentReq)
		require.NoError(t, err)
		assert.Nil(t, err)
		//act
		respStudentGet, err := studentRepo.GetByID(ctx, respStudent)
		//assert
		require.NoError(t, err)
		assert.Equal(t, testStudentReq.StudentName, respStudentGet.StudentName)
		assert.Equal(t, testStudentReq.Grade, respStudentGet.Grade)
	})
	t.Run("Fail Rows not effected", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)
		//act
		nonExistentStudentID := -1
		err = studentRepo.Update(ctx, int64(nonExistentStudentID), testStudentReq)
		require.Error(t, err)
		assert.Equal(t, pkgErrors.ErrNotFound, err)
	})
}
func TestDeleteStudent(t *testing.T) {

	db := postgres.NewFromEnv()
	defer db.DB.GetPool(context.Background()).Close()
	var (
		ctx           = context.Background()
		migrationPath = "./migrations"
	)
	t.Run("Success", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		respStudent, err := studentRepo.Add(ctx, testStudentReq)
		//assert
		require.NoError(t, err)
		assert.NotZero(t, respStudent)

		err = studentRepo.Delete(ctx, respStudent)
		require.NoError(t, err)
		assert.Nil(t, err)
		//act
		respStudentGet, err := studentRepo.GetByID(ctx, testStudentReq.StudentID)
		//assert
		require.Error(t, err)
		assert.Equal(t, serviceEntities.StudentRequest{}, respStudentGet)
	})
	t.Run("Failed to Delete, Not Found", func(t *testing.T) {
		db.SetUpDatabase(migrationPath)
		defer db.TearDownDatabase(migrationPath)
		//arrange
		studentRepo := NewStudentStorage(db.DB)
		testStudentReq := serviceEntities.StudentRequest{
			StudentID:   1,
			StudentName: "Test",
			Grade:       90,
		}
		//act
		err := studentRepo.Delete(ctx, testStudentReq.StudentID)
		require.Error(t, err)
		assert.ErrorIs(t, err, pkgErrors.ErrNotFound)
	})
}
