package main

import (
	"context"
	"fmt"
	"os"
	"post-master/database"
	"post-master/scraping"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"

	"cloud.google.com/go/storage"
)

const (
	// IsLocal 保存場所をローカルとGCSで分岐
	IsLocal = true

	TermMinutes = 1
	BaseRawURI  = "https://bun.uptrace.dev/guide/query-update.html"
	PostToDir   = "./post/"

	DBPATH = "./db.sqlite3"

	// Google Schedule バケツ名
	GCSBUCKETNAME = "hugo-storage-xxx.com"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	client := scraping.New(TermMinutes, BaseRawURI, &Data{})
	go client.Start(ctx, F)
	defer client.Close()

	// テスト用: プログラムをいつまで実行をするか
	time.Sleep(3 * time.Minute)
	cancel()
	time.Sleep(2 * time.Minute)
	log.Info().Msg("正しく終了")
}

// Data 対象ごとに欲しているデータ
type Data struct {
	ID string

	Name    string
	Tel     string
	Mail    string
	PostID  string
	Address string

	PostedAt  time.Time
	CreatedAt time.Time
}

func F(page playwright.Page) (string, error) {
	// 格納用データ
	data := &Data{}

	db := database.New(DBPATH, data)
	defer db.Close()

	// 情報の取捨選択
	title, err := page.Title()
	if err != nil {
		return "", err
	}

	data.Name = title

	titles, err := page.Locator("h2").AllInnerTexts()
	if err != nil {
		return "", err
	}

	data.Address = strings.Join(titles, ", ")

	// DBに保存
	if err := db.Create(data); err != nil {
		return "", err
	}

	if IsLocal {
		// ファイルに保存
		if err := toFile(data); err != nil {
			return "", err
		}
	} else {
		// ファイルに保存
		if err := saveFileToGCS(data); err != nil {
			return "", err
		}
	}

	links, err := page.Locator("a").All()
	if err != nil {
		return "", err
	}

	for _, v := range links {
		link, _ := v.InnerText()
		fmt.Println(link)
	}

	// NextTargetURI

	// save
	if err := db.Create(data); err != nil {
		return "", err
	}

	return "", nil
}

func toString(data *Data) string {
	// HUGOなどmarkdownで汎化性能が高そうな形式に整形してファイル保存する
	return fmt.Sprintf(
		`+++
date = %s
title = "%s"
description = ""
categories = []
tag = []
author = "" 
mainImage = ""
+++

## %s

名称
: %s
電話番号
: %s
メールアドレス
: %s
所在地
: %s
制作日
: %s


`,
		time.Now().String(),
		data.Name, // タイトル
		data.Name, // 記事本文

		// 構造化されたデータ
		data.Name,
		data.Tel,
		data.Mail,
		data.Address,
		data.CreatedAt.Format("2006年01月02日 15:04"),
	)
}

func toFile(data *Data) error {
	uuid := uuid.New()
	f, err := os.Create(fmt.Sprintf("%s%s.md", PostToDir, uuid.String()))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(toString(data))

	return nil
}

func saveFileToGCS(data *Data) error {
	// GCSクライアントの初期化
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal().Msgf("Failed to create GCS client: %v", err)
		return err
	}
	defer client.Close()

	// 保存するファイルの名前
	uuid := uuid.New()
	destFileName := fmt.Sprintf("%s%s.md", PostToDir, uuid.String())

	// GCSバケットへファイルを保存
	bucket := client.Bucket(GCSBUCKETNAME)
	obj := bucket.Object(destFileName)
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	if _, err := writer.Write([]byte(toString(data))); err != nil {
		log.Info().Msgf("Failed to write data to GCS: %v", err)
		return err
	}

	fmt.Printf("File %s uploaded to GCS bucket %s\n", destFileName, GCSBUCKETNAME)

	return nil
}
