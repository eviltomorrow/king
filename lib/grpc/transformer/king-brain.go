package transformer

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/mathutil"
)

func GenerateMarketStatusToMap(status *pb.MarketStatus) map[string]interface{} {
	return map[string]interface{}{
		"date":              status.Date,
		"shang_zheng_value": status.MarketIndex.ShangZheng.Value,
		"shang_zheng_direction": func() string {
			if status.MarketIndex.ShangZheng.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"shang_zheng_change": status.MarketIndex.ShangZheng.HasChanged,

		"shen_zheng_value": status.MarketIndex.ShenZheng.Value,
		"shen_zheng_direction": func() string {
			if status.MarketIndex.ShenZheng.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"shen_zheng_change": status.MarketIndex.ShenZheng.HasChanged,

		"chuang_ye_value": status.MarketIndex.ChuangYe.Value,
		"chuang_ye_direction": func() string {
			if status.MarketIndex.ChuangYe.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"chuang_ye_change": status.MarketIndex.ChuangYe.HasChanged,

		"ke_chuang50_value": status.MarketIndex.KeChuang_50.Value,
		"ke_chuang50_direction": func() string {
			if status.MarketIndex.KeChuang_50.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"ke_chuang50_change": status.MarketIndex.KeChuang_50.HasChanged,

		"bei_zheng50_value": status.MarketIndex.BeiZheng_50.Value,
		"bei_zheng50_direction": func() string {
			if status.MarketIndex.BeiZheng_50.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"bei_zheng50_change": status.MarketIndex.BeiZheng_50.HasChanged,

		"total":           status.MarketStockCount.Total,
		"rise_gt_7":       status.MarketStockCount.RiseGt_7,
		"rise_gt_7_ratio": mathutil.Trunc2(float64(status.MarketStockCount.RiseGt_7) / float64(status.MarketStockCount.Total) * 100),
		"fell_gt_7":       status.MarketStockCount.FellGt_7,
		"fell_gt_7_ratio": mathutil.Trunc2(float64(status.MarketStockCount.FellGt_7) / float64(status.MarketStockCount.Total) * 100),
	}
}
