package cache

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type LruCache struct {
	queue List
	items sync.Map

	maxSize int
	dir     string
}

type cacheItem struct {
	filename string
	key      string
}

func New(maxSize int, cacheDir string) *LruCache {
	cache := &LruCache{maxSize: maxSize, queue: NewList(), dir: cacheDir}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		dirErr := os.MkdirAll(cacheDir, os.ModePerm)
		if dirErr != nil {
			log.Panic("Cache dir unavaildable")
		}
	} else {
		// cache dir exists - should load cached files to mem
		go func() {
			cache.loadFilesCache()
		}()
	}

	return cache
}

func (c *LruCache) Get(key string) ([]byte, error) {
	cachedElem, inCache := c.items.Load(key)

	if inCache {
		item := cachedElem.(*ListItem)

		filename := item.Value.(cacheItem).filename
		content, err := c.readFile(filename)
		if err != nil {
			c.removeItemAndFile(item)
			return nil, err
		}
		c.queue.MoveToFront(item)
		return content, nil
	}

	return nil, nil
}

func (c *LruCache) Set(key string, value []byte) error {
	cachedElem, inCache := c.items.Load(key)

	filename := filepath.Join(c.dir, base64.StdEncoding.EncodeToString([]byte(key)))

	if inCache {
		item := cachedElem.(*ListItem)

		if item.Value.(cacheItem).filename == filename {
			c.queue.MoveToFront(item)
			return nil
		}
	}

	// remove last item
	if c.maxSize == c.queue.Len() {
		c.removeItemAndFile(c.queue.Back())
	}

	// write file
	err := os.WriteFile(filename, value, os.ModePerm)
	if err != nil {
		return err
	}

	item := c.queue.PushFront(cacheItem{filename: filename, key: key})
	c.items.Store(key, item)
	return nil
}

func (c *LruCache) loadFilesCache() {
	files, _ := os.ReadDir(c.dir)
	for _, file := range files {
		filename := file.Name()
		cacheKey, _ := base64.StdEncoding.DecodeString(filename)

		item := c.queue.PushFront(cacheItem{filename: filepath.Join(c.dir, filename), key: string(cacheKey)})
		c.items.Store(string(cacheKey), item)
	}
}

func (c *LruCache) readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c *LruCache) removeItemAndFile(item *ListItem) {
	cacheItem := item.Value.(cacheItem)
	c.queue.Remove(item)
	c.items.Delete(cacheItem.key)

	os.Remove(cacheItem.filename)
}
