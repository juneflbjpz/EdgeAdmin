// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package servergrouputils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// InitGroup 初始化分组信息
func InitGroup(parent *actionutils.ParentAction, groupId int64, menuItem string) (*pb.ServerGroup, error) {
	groupResp, err := parent.RPC().ServerGroupRPC().FindEnabledServerGroup(parent.AdminContext(), &pb.FindEnabledServerGroupRequest{ServerGroupId: groupId})
	if err != nil {
		return nil, err
	}
	var group = groupResp.ServerGroup
	if group == nil {
		return nil, errors.New("group with id '" + types.String(groupId) + "' not found")
	}

	parent.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	// 初始化设置菜单
	if len(menuItem) > 0 {
		// 获取设置概要信息
		configInfoResp, err := parent.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(parent.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{ServerGroupId: groupId})
		if err != nil {
			return group, err
		}

		var urlPrefix = "/servers/groups/group/settings"
		var leftMenuItems = []maps.Map{
			{
				"name":     "HTTP代理",
				"url":      urlPrefix + "/httpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "httpReverseProxy",
				"isOn":     configInfoResp.HasHTTPReverseProxy,
			},
			{
				"name":     "TCP代理",
				"url":      urlPrefix + "/tcpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "tcpReverseProxy",
				"isOn":     configInfoResp.HasTCPReverseProxy,
			},
			{
				"name":     "UDP代理",
				"url":      urlPrefix + "/udpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "udpReverseProxy",
				"isOn":     configInfoResp.HasUDPReverseProxy,
			},
		}

		leftMenuItems = filterMenuItems(leftMenuItems, groupId, urlPrefix, menuItem, configInfoResp)

		leftMenuItems = append(leftMenuItems, maps.Map{
			"name": "-",
			"url":  "",
		})
		leftMenuItems = append(leftMenuItems, maps.Map{
			"name":     "访客IP地址",
			"url":      urlPrefix + "/remoteAddr?groupId=" + types.String(groupId),
			"isActive": menuItem == "remoteAddr",
			"isOn":     configInfoResp.HasRemoteAddrConfig,
		})

		leftMenuItems = filterMenuItems2(leftMenuItems, groupId, urlPrefix, menuItem, configInfoResp)

		parent.Data["leftMenuItems"] = leftMenuItems
	}

	return group, nil
}
