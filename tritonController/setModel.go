package tritonController

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/juncheoll/triton-torrent/modelStoreCommunicator"
	"github.com/juncheoll/triton-torrent/setting"
	"github.com/juncheoll/triton-torrent/src/logCtrlr"
)

func SetModel(provider string, model string, version string, channel *chan string) error {
	filePath := fmt.Sprintf("%s/%s", setting.ModelsPath, provider)
	fileName := fmt.Sprintf("%s@%s#%s.torrent", provider, model, version)

	// Creating the provider folder.
	// If the provider folder already exists, it will not be created.
	logCtrlr.Log("Create provider folder.")
	if err := makeFolder(filePath); err != nil {
		return err
	}

	// downloading the model from the Model Store.
	logCtrlr.Log("Request model download to the Model Store.")
	modelFile, err := modelStoreCommunicator.GetModel(provider, model, version)
	if err != nil {
		return err
	}
	logCtrlr.Log("Successfully completed the torrent file download.")

	// Create a directory corresponding to the path.
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create and write to a file.
	file, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		return err
	}
	file.Write(modelFile)

	// Workdirectory 이동
	if err := os.Chdir(filePath); err != nil {
		log.Println(err)
		return nil
	}

	/*--------------------ModelStore 한테 토렌트 파일 받아와서 저장 완료--------------*/

	// File Download by Torrent
	// 토렌트 클라이언트 설정
	cfg := torrent.NewDefaultClientConfig()
	//cfg.Seed = true // 시딩 활성화
	cfg.ListenPort = 9000

	// 토렌트 클라이언트 생성
	cl, err := torrent.NewClient(cfg)
	if err != nil {
		log.Printf("error creating client: %s", err)
	}
	defer cl.Close()

	// 토렌트 파일 추가
	torrentPath := filePath + "/" + fileName // 여기에 토렌트 파일 경로 입력
	pwd, _ := os.Getwd()
	log.Println(pwd)
	metaInfo, err := metainfo.LoadFromFile(torrentPath)
	if err != nil {
		log.Printf("error loading torrent file: %s", err)
	}
	t, err := cl.AddTorrent(metaInfo)
	if err != nil {
		log.Fatalf("error adding torrent: %s", err)
	}

	<-t.GotInfo() // 토렌트 정보를 받을 때까지 대기

	t.DownloadAll()

	// 파일 다운로드 대기
	for !t.Complete.Bool() {
		time.Sleep(time.Second * 1)
		peers := t.PeerConns()
		if len(peers) != 0 {
			for _, peer := range peers {
				parts := strings.Split(peer.RemoteAddr.String(), ":")
				log.Println("Peer:", parts[0])
			}
		}
	}

	log.Printf("Downloaded %s", t.Name())
	*channel <- fileName

	/*
		if err := os.Chdir("../"); err != nil {
			log.Println(err)
			return nil
		}
	*/

	return nil
}
