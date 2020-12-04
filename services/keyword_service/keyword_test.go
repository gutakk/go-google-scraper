package keyword_service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

type KeywordServiceTestSuite struct {
	suite.Suite
	keywordService Keyword
}

func (s *KeywordServiceTestSuite) SetupTest() {
	s.keywordService = Keyword{}
}

func TestKeywordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordServiceTestSuite))
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithValidFileType() {
	result := s.keywordService.ValidateFileType("text/csv")

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithInvalidFileType() {
	result := s.keywordService.ValidateFileType("test")

	assert.Equal(s.T(), "File must be CSV format", result.Error())
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithMinRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1)

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithMaxRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1000)

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithZeroRow() {
	result := s.keywordService.ValidateCSVLength(0)

	assert.Equal(s.T(), "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithGreaterThanMaxRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1001)

	assert.Equal(s.T(), "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func (s *KeywordServiceTestSuite) TestReadFileWithValidFile() {
	result, err := s.keywordService.ReadFile("../../tests/fixture/adword_keywords.csv")

	assert.Equal(s.T(), []string{"AWS"}, result)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceTestSuite) TestReadFileWithFileNotFound() {
	result, err := s.keywordService.ReadFile("")

	assert.Equal(s.T(), nil, result)
	assert.Equal(s.T(), "Something went wrong, please try again", err.Error())
}
