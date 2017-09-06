package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func getMetricsStatistics(svcCloudWatch *cloudwatch.CloudWatch, startTime, endTime time.Time, nameSpace, metricsName *string, statistics string, demensions []*cloudwatch.Dimension) []float64 {
	respGetMetricStatistics, err := svcCloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  nameSpace,
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		MetricName: metricsName,
		Period:     aws.Int64(86400),
		Statistics: []*string{
			aws.String(statistics),
		},
		Dimensions: demensions,
	})

	stats := make([]float64, 0)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		jsonBody, _ := json.Marshal(respGetMetricStatistics)

		var result result
		json.Unmarshal(jsonBody, &result)
		sort.Sort(result.Datapoints)

		if len(result.Datapoints) > 0 {
			for _, dp := range result.Datapoints {
				switch statistics {
				case "Sum":
					stats = append(stats, dp.Sum)
				case "Maximum":
					stats = append(stats, dp.Maximum)
				case "Minimum":
					stats = append(stats, dp.Minimum)
				case "Average":
					stats = append(stats, dp.Average)
				case "SampleCount":
					stats = append(stats, dp.SampleCount)
				}
			}
			return stats
		}
	}
	return []float64{0}
}

func formatStorage(bytes float64) string {
	if bytes >= 1024*1024*1024*1024 {
		return fmt.Sprintf("%0.2f TB", bytes/(1024*1024*1024*1024))
	} else if bytes > 1024*1024*1024 {
		return fmt.Sprintf("%0.2f GB", bytes/(1024*1024*1024))
	} else if bytes > 1024*1024 {
		return fmt.Sprintf("%0.2f MB", bytes/(1024*1024))
	} else if bytes > 1024 {
		return fmt.Sprintf("%0.2f KB", bytes/(1024))
	}
	return fmt.Sprintf("%0.2f Bytes", bytes)
}
