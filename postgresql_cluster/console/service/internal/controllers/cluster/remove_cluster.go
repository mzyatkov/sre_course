package cluster

import (
	"encoding/base64"
	"encoding/json"
	"postgresql-cluster-console/internal/configuration"
	"postgresql-cluster-console/internal/controllers"
	"postgresql-cluster-console/internal/storage"
	"postgresql-cluster-console/internal/watcher"
	"postgresql-cluster-console/internal/xdocker"
	"postgresql-cluster-console/pkg/tracer"
	"postgresql-cluster-console/restapi/operations/cluster"

	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
)

type removeClusterHandler struct {
	db            storage.IStorage
	dockerManager xdocker.IManager
	logCollector  watcher.LogCollector
	log           zerolog.Logger
	cfg           *configuration.Config
}

func NewRemoveClusterHandler(db storage.IStorage, dockerManager xdocker.IManager, logCollector watcher.LogCollector, cfg *configuration.Config, log zerolog.Logger) cluster.PostClustersIDRemoveHandler {
	return &removeClusterHandler{
		db:            db,
		dockerManager: dockerManager,
		logCollector:  logCollector,
		log:           log,
		cfg:           cfg,
	}
}

func (h *removeClusterHandler) Handle(param cluster.PostClustersIDRemoveParams) middleware.Responder {
	cid := param.HTTPRequest.Context().Value(tracer.CtxCidKey{}).(string)
	localLog := h.log.With().Str("cid", cid).Logger()
	clusterInfo, err := h.db.GetCluster(param.HTTPRequest.Context(), param.ID)
	if err != nil {
		return cluster.NewPostClustersIDRemoveBadRequest().WithPayload(controllers.MakeErrorPayload(err, controllers.BaseError))
	}

	var extraVars []string

	err = json.Unmarshal(clusterInfo.ExtraVars, &extraVars)
	if err != nil {
		return cluster.NewPostClustersIDRemoveBadRequest().WithPayload(controllers.MakeErrorPayload(err, controllers.BaseError))
	}
	extraVars = append(extraVars, "state=absent")

	var (
		envs          []string
		paramLocation ParamLocation
	)
	if clusterInfo.SecretID != nil {
		envs, paramLocation, err = getSecretEnvs(param.HTTPRequest.Context(), h.log, h.db, *clusterInfo.SecretID, h.cfg.EncryptionKey)
		if err != nil {
			return cluster.NewPostClustersIDRemoveBadRequest().WithPayload(controllers.MakeErrorPayload(err, controllers.BaseError))
		}
		if paramLocation == ExtraVarsParamLocation {
			extraVars = append(extraVars, envs...)
		}
	}
	envs = append(envs, "patroni_cluster_name="+clusterInfo.Name)
	if len(clusterInfo.Inventory) != 0 {
		envs = append(envs, "ANSIBLE_INVENTORY_JSON="+base64.StdEncoding.EncodeToString(clusterInfo.Inventory))
	}
	localLog.Trace().Strs("envs", envs).Msg("got envs")

	dockerId, err := h.dockerManager.ManageCluster(param.HTTPRequest.Context(), &xdocker.ManageClusterConfig{
		Envs:      envs,
		ExtraVars: extraVars,
	})
	if err != nil {
		return cluster.NewPostClustersIDRemoveBadRequest().WithPayload(controllers.MakeErrorPayload(err, controllers.BaseError))
	}
	localLog.Trace().Str("docker_id", string(dockerId)).Msg("docker was started")

	err = h.db.DeleteCluster(param.HTTPRequest.Context(), clusterInfo.ID)
	if err != nil {
		return cluster.NewPostClustersIDRemoveBadRequest().WithPayload(controllers.MakeErrorPayload(err, controllers.BaseError))
	}

	h.logCollector.PrintToConsole(dockerId, cid)

	return cluster.NewPostClustersIDRemoveNoContent()
}
