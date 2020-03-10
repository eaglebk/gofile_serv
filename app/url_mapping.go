package app

import "space.eagle2000/fileserv/controllers/files"

func mapUrls() {
	router.GET("/files/:filename", files.GetFile)
	router.POST("/files/upload", files.UploadFile)
	router.PUT("/files/update", files.UpdateFile)
	router.DELETE("/files/delete/:filename", files.DeleteFile)
}
