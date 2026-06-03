package management

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/managementasset"
)

// RefreshManagementAsset triggers a management control panel asset sync immediately.
func (h *Handler) RefreshManagementAsset(c *gin.Context) {
	if h == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "handler_unavailable", "message": "management handler is unavailable"})
		return
	}

	h.mu.Lock()
	cfg := h.cfg
	configFilePath := h.configFilePath
	h.mu.Unlock()

	proxyURL := ""
	panelRepository := ""
	if cfg != nil {
		if cfg.RemoteManagement.DisableControlPanel {
			c.JSON(http.StatusConflict, gin.H{"error": "control_panel_disabled", "message": "management control panel is disabled"})
			return
		}
		proxyURL = strings.TrimSpace(cfg.ProxyURL)
		panelRepository = strings.TrimSpace(cfg.RemoteManagement.PanelGitHubRepository)
	}

	staticDir := managementasset.StaticDir(configFilePath)
	if staticDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "static_dir_unavailable", "message": "management static directory is unavailable"})
		return
	}

	exists, err := managementasset.RefreshLatestManagementHTML(c.Request.Context(), staticDir, proxyURL, panelRepository)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "refresh_failed", "message": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusBadGateway, gin.H{"error": "refresh_failed", "message": "failed to refresh management control panel asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
