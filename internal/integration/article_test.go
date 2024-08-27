package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"testing"
)

// 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
}

func (s *ArticleTestSuite) SetupSuite() {

}

func (s *ArticleTestSuite) TestEdit() {

}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
