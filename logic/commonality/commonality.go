package commonality

import (
	"Go-Live/global"
	receive "Go-Live/interaction/receive/commonality"
	response "Go-Live/interaction/response/commonality"
	"Go-Live/models/config/uploadMethod"
	"Go-Live/models/contribution/video"
	"Go-Live/models/users"
	"Go-Live/models/users/attention"
	"Go-Live/utils/conversion"
	"Go-Live/utils/location"
	"Go-Live/utils/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

var (
	//Temporary 文件文件存续位置
	Temporary = "assets/tmp/"
)

func OssSTS() (results interface{}, err error) {
	info, err := oss.GteStsInfo()
	if err != nil {
		return nil, err
	}
	res, err := response.GteStsInfo(info)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

func Upload(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//上传文件明
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("上传接口不存在")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("请联系管理员设置接口保存路径")
	}
	//取出文件
	index := strings.LastIndex(file.Filename, ".")
	suffix := file.Filename[index:]
	switch suffix {
	case ".jpg", ".jpeg", ".png", ".ico", ".gif", ".wbmp", ".bmp", ".svg", ".webp", ".mp4":
	default:
		return nil, fmt.Errorf("非法后缀！")
	}
	if !location.IsDir(method.Path) {
		if err = os.MkdirAll(method.Path, 077); err != nil {
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}
	dst := method.Path + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Warn("update headPortrait err")
		return nil, fmt.Errorf("上传失败")
	} else {
		return dst, nil
	}
}

func UploadSlice(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//上传文件明
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("上传接口不存在")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("请联系管理员设置接口保存路径")
	}
	if !location.IsDir(Temporary) {
		if err = os.MkdirAll(Temporary, 077); err != nil {
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}
	dst := Temporary + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Warn("Possible cause of the fragment upload failure: ", err)
		return nil, fmt.Errorf("上传失败")
	} else {
		return dst, nil
	}
}

func UploadCheck(data *receive.UploadCheckStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("未配置上传方法")
	}
	list := make(receive.UploadSliceList, 0)
	path := method.Path + "/" + data.FileMd5
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//文件已存在
		return response.UploadCheckResponse(true, list, path)
	}
	//取出未上传的分片
	for _, v := range data.SliceList {
		if _, err := os.Stat(Temporary + "/" + v.Hash); os.IsNotExist(err) {
			list = append(list, receive.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	return response.UploadCheckResponse(false, list, "")
}

func UploadMerge(data *receive.UploadMergeStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("未配置上传方法")
	}
	dst := method.Path + "/" + data.FileName
	list := make(receive.UploadSliceList, 0)
	path := method.Path + "/" + data.FileName
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//文件已存在直接返回
		return dst, nil
	}
	//取出未上传的分片
	for _, v := range data.SliceList {
		if _, err := os.Stat(Temporary + "/" + v.Hash); os.IsNotExist(err) {
			list = append(list, receive.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	if len(list) > 0 {
		return nil, fmt.Errorf("分片未全部上传")
	}
	//进行合并操作
	cf, err := os.Create(dst)
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			global.Logger.Errorf("合并操作释放内存失败 %d", err)
		}
	}(cf)
	if err != nil {
		return nil, fmt.Errorf("保存失败")
	}
	fileInfo, err := os.OpenFile(dst, os.O_APPEND, os.ModeSetuid)
	defer func(fileInfo *os.File) {
		if err := fileInfo.Close(); err != nil {
			global.Logger.Errorf("关闭资源 err : %d", err)
		}
	}(fileInfo)
	//合并操作
	for _, v := range data.SliceList {
		tmpFile, err := os.OpenFile(Temporary+"/"+v.Hash, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		b, err := ioutil.ReadAll(tmpFile)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := fileInfo.Write(b); err != nil {
			global.Logger.Errorf("合并分片追加错误 err : %d", err)
		}
		// 关闭分片
		if err := tmpFile.Close(); err != nil {
			global.Logger.Errorf("关闭分片错误 err : %d", err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			global.Logger.Errorf("合并操作删除临时分片失败 err : %d", err)
		}
	}
	return dst, nil
}

func UploadingMethod(data *receive.UploadingMethodStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if method.IsExistByField("interface", data.Method) {
		return response.UploadingMethodResponse(method.Method), nil
	} else {
		return nil, fmt.Errorf("未配置上传方法")
	}
}

func UploadingDir(data *receive.UploadingDirStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if method.IsExistByField("interface", data.Interface) {
		return response.UploadingDirResponse(method.Path), nil
	} else {
		return nil, fmt.Errorf("未配置上传方法")
	}
}

func GetFullPathOfImage(data *receive.GetFullPathOfImageMethodStruct) (results interface{}, err error) {
	path, err := conversion.SwitchIngStorageFun(data.Type, data.Path)
	if err != nil {
		return nil, err
	}
	return path, nil
}

func Search(data *receive.SearchStruct, uid uint) (results interface{}, err error) {
	switch data.Type {
	case "video":
		//视频搜索
		list := new(video.VideosContributionList)
		err = list.Search(data.PageInfo)
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
		res, err := response.SearchVideoResponse(list)
		if err != nil {
			return nil, fmt.Errorf("响应失败")
		}
		return res, nil
		break
	case "user":
		list := new(users.UserList)
		err := list.Search(data.PageInfo)
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
		aids := make([]uint, 0)
		if uid != 0 {
			//用户登入情况下
			al := new(attention.AttentionsList)
			err = al.GetAttentionList(uid)
			if err != nil {
				return nil, fmt.Errorf("获取关注列表失败")
			}
			for _, v := range *al {
				aids = append(aids, v.AttentionID)
			}
		}
		res, err := response.SearchUserResponse(list, aids)
		return res, nil
		break
	default:
		return nil, fmt.Errorf("未匹配的类型")
	}
	return
}
