package convert

import (
	"postgresql-cluster-console/internal/storage"
	"postgresql-cluster-console/models"

	"github.com/go-openapi/strfmt"
)

func ClusterToSwagger(cl *storage.Cluster, servers []storage.Server, environmentCode, projectCode string) *models.ClusterInfo {
	clusterInfo := &models.ClusterInfo{
		ConnectionInfo: cl.ConnectionInfo,
		CreationTime:   strfmt.DateTime(cl.CreatedAt),
		ClusterLocation: func() string {
			if cl.Location != nil {
				return *cl.Location
			}

			return ""
		}(),
		Environment:     environmentCode,
		ID:              cl.ID,
		Servers:         make([]*models.ClusterInfoInstance, 0, len(servers)),
		Name:            cl.Name,
		Description:     cl.Description,
		PostgresVersion: cl.PostgreVersion,
		ProjectName:     projectCode,
		Status:          cl.Status,
	}

	for _, server := range servers {
		clusterInfo.Servers = append(clusterInfo.Servers, &models.ClusterInfoInstance{
			ID:             server.ID,
			IP:             server.IpAddress.String(),
			Lag:            server.Lag,
			Name:           server.Name,
			PendingRestart: server.PendingRestart,
			Role:           server.Role,
			Status:         server.Status,
			Tags:           server.Tags,
			Timeline:       server.Timeline,
		})
	}

	return clusterInfo
}
