package main

import (
	"context"
	"time"

	"github.com/TarsDemo/Tars-MiniProgramm-Service-ClubActivityServer/tars-protocol/LifeService"
)

// ClubActivityManagerImp servant implementation
type ClubActivityManagerImp struct {
	dataServiceProxy     *LifeService.DataService
	dataServiceObj       string
	userInfoServiceProxy *LifeService.UserInfoService
	userInfoServiceObj   string
}

// Init servant init
func (imp *ClubActivityManagerImp) Init() error {
	//initialize servant here:
	imp.dataServiceProxy = new(LifeService.DataService)
	imp.dataServiceObj = "LifeService.DataServer.DataServiceObj"
	imp.userInfoServiceProxy = new(LifeService.UserInfoService)
	imp.userInfoServiceObj = "LifeService.UserInfoServer.UserInfoServiceObj"

	comm.StringToProxy(imp.dataServiceObj, imp.dataServiceProxy)
	comm.StringToProxy(imp.userInfoServiceObj, imp.userInfoServiceProxy)
	return nil
}

// Destroy servant destory
func (imp *ClubActivityManagerImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *ClubActivityManagerImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}

//CreateClubManager 创建社团管理员
func (imp *ClubActivityManagerImp) CreateClubManager(ctx context.Context, wxID string, clubID string, errCode *LifeService.ErrorCode) (int32, error) {
	// 判断是否已存在管理员
	var isClubManager bool
	_, err1 := imp.userInfoServiceProxy.IsClubManager(wxID, clubID, &isClubManager)
	if err1 != nil {
		SLOG.Error("Call Remote UserInfoServer::IsClubManager error: ", err1.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	if isClubManager {
		*errCode = LifeService.ErrorCode_MANAGEREXIST
		return 0, nil
	}
	// 不存在则新建管理员
	_, err := imp.dataServiceProxy.CreateClubManager(wxID, clubID)
	if err != nil {
		SLOG.Error("Call Remote DataServer::CreateClubManager error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, err
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("CreateClubManager successfully")
	return 0, nil
}

//CreateClub 创建社团
func (imp *ClubActivityManagerImp) CreateClub(ctx context.Context, clubInfo *LifeService.ClubInfo, errCode *LifeService.ErrorCode) (int32, error) {
	var clubID string
	// 获取当前时间
	CurrentTime := time.Now().Format("2006-01-02 15:04:05")
	clubInfo.Create_time = CurrentTime
	// 创建社团
	_, err := imp.dataServiceProxy.CreateClub(clubInfo, &clubID)

	if err != nil {
		SLOG.Error("Create Club Error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return -1, err
	}
	// 返回的社团id为空, 说明DB层服务SQL执行错误
	if clubID == "" {
		SLOG.Error("Create Club Error")
		*errCode = LifeService.ErrorCode_CLUBEXIST
		return 0, nil
	}
	// 创建社团管理员
	_, err1 := imp.dataServiceProxy.CreateClubManager(clubInfo.Chairman, clubID)
	if err1 != nil {
		SLOG.Error("Call Remote DataServer::CreateClubManager error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("CreateClubManager successfully")
	return 0, nil
}

//GetClubList 获取社团列表
func (imp *ClubActivityManagerImp) GetClubList(ctx context.Context, index int32, wxID string, nextIndex *int32, clubInfoList *[]LifeService.ClubInfo, errCode *LifeService.ErrorCode) (int32, error) {
	// 一次获取数量为6
	var batch int32 = 6
	iRet, err := imp.dataServiceProxy.GetClubList(index, batch, wxID, nextIndex, clubInfoList)
	if err != nil {
		SLOG.Error("Get club list error")
		*errCode = LifeService.ErrorCode_SERVERERROR
	} else {
		if iRet == 0 {
			SLOG.Debug("GetClubList successfully")
			*errCode = LifeService.ErrorCode_SUCCESS
		} else {
			SLOG.Debug("Cannot get Club list")
			*errCode = LifeService.ErrorCode_SERVERERROR
		}
	}
	return 0, nil
}

//GetManagerClubList 获取管理社团列表
func (imp *ClubActivityManagerImp) GetManagerClubList(ctx context.Context, index int32, wxID string, nextIndex *int32, clubInfoList *[]LifeService.ClubInfo, errCode *LifeService.ErrorCode) (int32, error) {
	var batch int32 = 6

	iRet, err := imp.dataServiceProxy.GetManagerClubList(index, batch, wxID, nextIndex, clubInfoList)
	if err != nil {
		SLOG.Error("Get club list error")
		*errCode = LifeService.ErrorCode_SERVERERROR
	} else {
		if iRet == 0 {
			SLOG.Debug("GetClubList successfully")
			*errCode = LifeService.ErrorCode_SUCCESS
		} else {
			SLOG.Debug("Cannot get Club list")
			*errCode = LifeService.ErrorCode_SERVERERROR
		}
	}
	return 0, nil
}

//DeleteClub 删除社团
func (imp *ClubActivityManagerImp) DeleteClub(ctx context.Context, clubID string, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.DeleteClub(clubID, &affectRows)
	if err != nil || affectRows == -1 {
		SLOG.Error("Remote Server DataServer::DeleteClub error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if affectRows == 0 {
		*errCode = LifeService.ErrorCode_CLUBNOTEXIST
		SLOG.Error("Remote Server DataServer::DeleteClub error: Club does not exist")
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("Delete Club Successfully")
	return 0, nil
}

//DeleteClubManager 删除社团管理员
func (imp *ClubActivityManagerImp) DeleteClubManager(ctx context.Context, wxID string, clubID string, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.DeleteClubManager(wxID, clubID, &affectRows)
	if err != nil || affectRows == -1 {
		SLOG.Error("Remote Server DataServer::DeleteClubManager Error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if affectRows == 0 {
		*errCode = LifeService.ErrorCode_MANAGERNOTEXIST
		SLOG.Error("Remote Server DataServer::DeleteClubManager error: Manager does not exist")
		return 0, nil
	}

	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("Delete Club Manager Successfully")
	return 0, nil
}

//ApplyForClub 申请社团
func (imp *ClubActivityManagerImp) ApplyForClub(ctx context.Context, wxID string, clubID string, errCode *LifeService.ErrorCode) (int32, error) {
	// 判断是否已经在社团中或提交了申请
	var isInClub bool
	_, err := imp.userInfoServiceProxy.IsInClub(wxID, clubID, false, &isInClub)

	if err != nil {
		SLOG.Error("Remote Server UserInfoServer::IsInClub error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if isInClub {
		*errCode = LifeService.ErrorCode_USERAPPLIED
		SLOG.Debug("Applied")
		return 0, nil
	}

	// 若没有提交申请，创建申请
	_, err1 := imp.dataServiceProxy.CreateApply(wxID, clubID)
	if err1 != nil {
		SLOG.Error("Call Remote DataServer::CreateApply error", err1)
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("Apply Successfully")
	return 0, nil
}

//GetClubApply 获取社团申请列表
func (imp *ClubActivityManagerImp) GetClubApply(ctx context.Context, clubID string, index int32, applyStatus int32, nextIndex *int32, applyList *[]LifeService.ApplyInfo) (int32, error) {
	var batch int32 = 6
	iRet, err := imp.dataServiceProxy.GetApplyListByClubId(clubID, index, batch, applyStatus, nextIndex, applyList)
	if err != nil {
		SLOG.Error("Remote Server DataServer::GetApplyListByClubId error: ", err.Error())
		return -1, err
	}
	return iRet, err
}

//GetUserApply 获取用户的申请
func (imp *ClubActivityManagerImp) GetUserApply(ctx context.Context, wxID string, index int32, applyStatus int32, nextIndex *int32, applyList *[]LifeService.ApplyInfo) (int32, error) {
	var batch int32 = 6
	iRet, err := imp.dataServiceProxy.GetApplyListByUserId(wxID, index, batch, applyStatus, nextIndex, applyList)
	if err != nil {
		SLOG.Error("Remote Server DataServer::GetApplyListByUserId error: ", err.Error())
		return -1, err
	}
	return iRet, err
}

//ModifyApplyStatus 设置申请状态
func (imp *ClubActivityManagerImp) ModifyApplyStatus(ctx context.Context, wxID string, clubID string, applyStatus int32, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.SetApplyStatus(wxID, clubID, applyStatus, &affectRows)
	if err != nil || affectRows == -1 {
		SLOG.Error("Remote Server DataServer::setApplyStatus error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	// 没有改变任何记录，说明申请不存在
	if affectRows == 0 {
		SLOG.Error("Remote Server DataServer::setApplyStatus error: No such apply")
		*errCode = LifeService.ErrorCode_APPLYNOTEXIST
		return 0, nil
	}
	SLOG.Debug("ModifuApplyStatus")
	*errCode = LifeService.ErrorCode_SUCCESS
	return 0, nil
}

//DeleteApply 删除申请
func (imp *ClubActivityManagerImp) DeleteApply(ctx context.Context, wxID string, clubID string, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.DeleteApply(wxID, clubID, &affectRows)
	if err != nil {
		SLOG.Error("Remote Server DataServer::DeleteApply error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	// 没有改变任何记录，说明申请不存在
	if affectRows == 0 {
		*errCode = LifeService.ErrorCode_APPLYNOTEXIST
		SLOG.Error("Remote Server DataServer::DeleteApply error: Apply Not Exist")
		return 0, nil
	}

	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("DeleteApply")
	return 0, nil
}

//CreateActivity 创建活动
func (imp *ClubActivityManagerImp) CreateActivity(ctx context.Context, wxID string, activityInfo *LifeService.ActivityInfo, errCode *LifeService.ErrorCode) (int32, error) {
	var isClubManager bool
	// 判断是否为社团管理员
	_, err := imp.userInfoServiceProxy.IsClubManager(wxID, activityInfo.Club_id, &isClubManager)

	if err != nil {
		SLOG.Error("Remote Server UserInfoServer::IsClubManager error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if isClubManager {
		// 创建活动
		_, err1 := imp.dataServiceProxy.CreateActivity(activityInfo)

		if err1 != nil {
			SLOG.Error("Cal Remote DataServer::CreateActivity error ", err1.Error())
			*errCode = LifeService.ErrorCode_SERVERERROR
			return 0, nil
		}
		SLOG.Debug("Create Activity successfully")
		*errCode = LifeService.ErrorCode_SUCCESS
	} else {
		*errCode = LifeService.ErrorCode_MANAGERNOTEXIST
		return 0, nil
	}
	return 0, nil
}

//GetActivityList 获取活动列表
func (imp *ClubActivityManagerImp) GetActivityList(ctx context.Context, index int32, wxID string, clubID string, nextIndex *int32, activityList *[]map[string]string) (int32, error) {
	var batch int32 = 6

	_, err := imp.dataServiceProxy.GetActivityList(index, batch, wxID, clubID, nextIndex, activityList)

	if err != nil {
		SLOG.Error("Call Remote DataServer::GetActivityList error: ", err.Error())
		return -1, err
	}
	return 0, nil
}

//UpdateActivity 更新活动信息
func (imp *ClubActivityManagerImp) UpdateActivity(ctx context.Context, activityInfo *LifeService.ActivityInfo, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.UpdateActivity(activityInfo, &affectRows)

	if err != nil {
		SLOG.Error("Call Remote DataServer::UpdateActivity error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if affectRows == 0 {
		SLOG.Error("Remote DateServer::UpdateActivity Execute SQL error")
		*errCode = LifeService.ErrorCode_ACTIVITYNOTEXIST
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	return 0, nil
}

//DeleteActivity 删除活动
func (imp *ClubActivityManagerImp) DeleteActivity(ctx context.Context, activityID string, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.DeleteActivity(activityID, &affectRows)
	if err != nil {
		SLOG.Error("Remote Server DataServer::deleteActivity error: ", err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	if affectRows == 0 {
		*errCode = LifeService.ErrorCode_ACTIVITYNOTEXIST
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	SLOG.Debug("DeleteActivity")
	return 0, nil
}

//GetActivityDetail 获取活动详情
func (imp *ClubActivityManagerImp) GetActivityDetail(ctx context.Context, activityID string, activityInfo *LifeService.ActivityInfo) (int32, error) {

	_, err := imp.dataServiceProxy.GetActivityInfo(activityID, activityInfo)

	if err != nil {
		SLOG.Error("Call Remote DataServer::createActivity error: ", err.Error())
		return -1, err
	}

	return 0, nil
}

//GetActivityParticipate 获取活动参与者信息
func (imp *ClubActivityManagerImp) GetActivityParticipate(ctx context.Context, index int32, activityID string, nextIndex *int32, participateList *[]LifeService.ActivityRecord) (int32, error) {
	var batch int32 = 6
	iRet, err := imp.dataServiceProxy.GetActivityRecords(index, batch, activityID, nextIndex, participateList)

	if err != nil || iRet != 0 {
		SLOG.Error("Call Remote DataServer::getActivityRecords error")
		return -1, nil
	}
	return 0, nil
}

//ApplyForActivity 活动报名
func (imp *ClubActivityManagerImp) ApplyForActivity(ctx context.Context, wxID string, activityID string, errCode *LifeService.ErrorCode) (int32, error) {
	var isApplied bool
	// 是否已经报名活动
	_, err := imp.userInfoServiceProxy.IsAppliedActivity(wxID, activityID, &isApplied)
	if err != nil {
		SLOG.Error("Remote Server UserInfoServer::IsApplied error")
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	if isApplied {
		SLOG.Debug("Applied")
		*errCode = LifeService.ErrorCode_RECORDEXIST
		return 0, nil
	}

	// 未报名则创建活动报名记录
	_, err1 := imp.dataServiceProxy.CreateActivityRecord(wxID, activityID)
	if err1 != nil {
		SLOG.Error("Call Remote DataServer::createActivityRecord error", err1.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}

	SLOG.Debug("Apply Activity successfully")
	*errCode = LifeService.ErrorCode_SUCCESS

	return 0, nil
}

//DeleteActivityParticipate 删除参与者信息
func (imp *ClubActivityManagerImp) DeleteActivityParticipate(ctx context.Context, activityID string, wxID string, errCode *LifeService.ErrorCode) (int32, error) {
	var affectRows int32
	_, err := imp.dataServiceProxy.DeleteActivityRecord(activityID, wxID, &affectRows)

	if err != nil {
		SLOG.Error("Call Remote DataServer::deleteActivityRecord error: " + err.Error())
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	if affectRows == 0 {
		*errCode = LifeService.ErrorCode_RECORDNOTEXIST
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_SUCCESS
	return 0, nil
}
