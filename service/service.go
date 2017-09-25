package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
)

type Service struct {
	logger *logrus.Logger
	wg     *sync.WaitGroup
	mu     *sync.Mutex
	cache  map[string]*entry
}

type Schema struct {
	Url      string                 `json:"url"`
	Meta     map[string]interface{} `json:"meta"`
	Elements map[string]int         `json:"elemets"`
}

type entry struct {
	schema    *Schema
	processed bool
	err       error
	ready     chan struct{}
}

func NewService(logger *logrus.Logger) *Service {
	return &Service{
		logger: logger,
		wg:     &sync.WaitGroup{},
		mu:     &sync.Mutex{},
		cache:  make(map[string]*entry),
	}
}

func (s *Service) GetSchema(urls []string) ([]*Schema, error) {
	ret := make([]*Schema, 0, len(urls))
	for _, u := range urls {
		s.wg.Add(1)
		go func(url string) {
			s.mu.Lock()
			e := s.cache[url]

			if e == nil {
				e = &entry{
					ready: make(chan struct{}),
				}
				s.cache[url] = e
				s.mu.Unlock()

				e.schema, e.err = s.processUrl(url)
				e.processed = true
				close(e.ready)
			} else {
				s.mu.Unlock()
				if e.processed {
					s.logger.Infof("Url already processed: %q. Skip", url)
				} else {
					s.logger.Infof("Url is processing: %q. Please wait..", url)
				}
				<-e.ready
			}

			if e.err != nil {
				s.logger.Errorf("Error while processing url %q: %v", url, e.err)
				ret = append(ret, &Schema{})
			} else {
				s.logger.Infof("Url processed successfully: %q", url)
				ret = append(ret, e.schema)
			}

			s.wg.Done()
		}(u)
	}
	s.wg.Wait()

	return ret, nil
}

func (s *Service) processUrl(url string) (*Schema, error) {
	start := time.Now()
	s.logger.Debugf("Start getting html content from url: %q", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while loading content: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while parsing content: %v", err)
	}
	s.logger.Debugf(
		"Finish getting html content from url: %q. Time elapsed: %q (%d bytes downloaded)",
		url,
		time.Since(start),
		len(body),
	)

	txt := string(body)
	doc, err := html.Parse(strings.NewReader(txt))
	if err != nil {
		s.logger.Fatal(err)
	}

	elements := make(map[string]int)
	s.countHtmlTags(doc, elements)
	meta := s.getMeta(resp)

	return &Schema{
		Url:      url,
		Meta:     meta,
		Elements: elements,
	}, nil
}

func (s *Service) getMeta(resp *http.Response) map[string]interface{} {
	return map[string]interface{}{
		"status":         resp.StatusCode,
		"content-type":   resp.Header.Get("Content-type"),
		"content-length": resp.ContentLength,
	}
}

func (s *Service) countHtmlTags(n *html.Node, counter map[string]int) {
	if n.Type == html.ElementNode {
		if _, ok := counter[n.Data]; ok {
			counter[n.Data]++
		} else {
			counter[n.Data] = 1
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.countHtmlTags(c, counter)
	}
}
