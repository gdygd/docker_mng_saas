package gapi

import (
	"fmt"
	"io"
	"log"

	"agent-service/internal/logger"
	"agent-service/pb"
)

func (s *Server) DataStream(stream pb.ContainerService_DataStreamServer) error {
	logger.Log.Print(1, "DataStream...!")
	req, err := stream.Recv()
	if err == io.EOF {
		log.Println("client closed stream")
		return err
	}
	if err != nil {
		log.Printf("receive error: %v", err)
		return err
	}
	log.Printf("From client: %v", req)

	msg := &pb.ServerMessage{
		Command:         pb.CommandType_ACK,
		TargetContainer: "",
		Host:            "1",
	}

	if err := stream.Send(msg); err != nil {
		log.Printf("send error: %v", err)
		return err
	}

	return nil
}

func (s *Server) ConnMessage(stream pb.ContainerService_ConnMessageServer) error {
	// 채널 생성
	clientMsgs := make(chan *pb.Hello)

	respChan := make(chan *pb.Hello, 100)

	// recv
	go func() {
		defer close(clientMsgs)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				log.Println("client closed stream")
				return
			}
			if err != nil {
				log.Printf("receive error: %v", err)
				return
			}
			log.Printf("From client: %s", req.Msg)
			clientMsgs <- req

			s.pushJob(req, respChan)

		}
	}()

	// send
	for {
		select {
		case msg, ok := <-clientMsgs:
			if !ok {
				log.Println("client message channel closed")
				return nil
			}

			resp := &pb.Hello{Msg: fmt.Sprintf("Hello Client, you said: %s", msg.Msg)}
			if err := stream.Send(resp); err != nil {
				log.Printf("send error: %v", err)
				return err
			}

			serverPush := &pb.Hello{Msg: "Hello Server"}
			if err := stream.Send(serverPush); err != nil {
				log.Printf("send error: %v", err)
				return err
			}

		case resp := <-respChan:
			logger.Log.Print(2, "respchan .. %v", resp)

		case <-stream.Context().Done():
			log.Println("client disconnected")
			return nil
		}
	}
}

func (s *Server) pushJob(req *pb.Hello, res chan *pb.Hello) {
	logger.Log.Print(2, "pushJob")
	job := Job_ConnMessage{
		req, res,
	}
	s.jobCh <- job
}

// shutdown에서 worker종료
// job 확인 후
func (s *Server) worker(id int) {
	logger.Log.Print(2, "Worker %d start", id)
	defer func() {
		s.work_wg.Done()
		logger.Log.Print(2, "Ended worker %d", id)
	}()
	logger.Log.Print(2, "start work (%d)", id)

	for {
		select {
		case job, ok := <-s.jobCh:
			if !ok {
				logger.Log.Print(2, "worker %d: job closed..", id)
				return
			}

			logger.Log.Print(2, "woker : (%d) job : %v", id, job.Req)
			job.RespChan <- &pb.Hello{Msg: "from worker, Hello Server"}
		case <-s.work_ctx.Done():
			logger.Log.Print(2, "worker (%d) is stopping..", id)
			return
		}
	}
}
