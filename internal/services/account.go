package services

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/starjun/gotools"
	"log"
	"toes/internal/models"
)

var Account = new(accountService)

type accountService struct {
}

func (srv *accountService) FilterQueryFromResult(c context.Context, _reqParam *models.QueryConfigRequest) (ret []models.Account, totalCount int, err error) {
	reqMap := make(map[string]interface{})
	reqMap["offset"] = 0
	reqMap["limit"] = 500
	resp, cnt, err := models.AccountQueryList(c, reqMap)
	if cnt < 1 {
		return ret, totalCount, err
	}
	// 将contains转化成in
	for k, rule := range _reqParam.Query {
		if rule.Opt == models.ContainOpt {
			_reqParam.Query[k].Opt = models.InOpt
		}
	}
	var gotoolsRule []gotools.CRule
	err = mapstructure.Decode(_reqParam.Query, &gotoolsRule)
	if err != nil {
		log.Println("_reqParam.Query decode err", err)
	}
	for _, account := range resp {
		tmpMap := make(map[string]string)
		err = mapstructure.Decode(account, &tmpMap)
		if err != nil {
			log.Println("account decode", err)
			continue
		}
		isFilterPass := gotools.MapCrulesListMatch(tmpMap, gotoolsRule)
		if !isFilterPass {
			continue
		}
		totalCount++
		ret = append(ret, account)
	}
	limit := _reqParam.Limit
	if limit > 100 {
		limit = 100
	}
	stop := _reqParam.Offset + limit
	totalCount = len(ret)
	if totalCount < 1 {
		return ret, totalCount, err
	}
	if _reqParam.Offset > totalCount {
		_reqParam.Offset = totalCount
	}
	if totalCount < stop {
		stop = totalCount
	}
	ret = ret[_reqParam.Offset:stop]

	return ret, totalCount, err
}
