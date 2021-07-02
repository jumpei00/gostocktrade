package models_test

import (
	"github.com/jumpei00/gostocktrade/app/models"
)

func (suite *ModelsTestSuite) TestCreateBacktestResult() {
	suite.Nil(suite.Op.CreateBacktestResult())

	models.DeleteBacktestResult("VOO")
}

func (suite *ModelsTestSuite) TestGetOptimizedParamFrame() {
	// initializing
	suite.Op.CreateBacktestResult()

	opframe := models.GetOptimizedParamFrame("VOO")
	suite.NotEmpty(opframe.Param)

	opframe = models.GetOptimizedParamFrame("TEST")
	suite.Nil(opframe.Param)

	models.DeleteBacktestResult("VOO")
	opframe = models.GetOptimizedParamFrame("VOO")
	suite.Nil(opframe.Param)
}
