package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/distribution/context"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/joho/godotenv"
)

var cli *client.Client

func main() {
	fmt.Println("Helow")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env File Founf ", err)
	}

	val, _ := os.LookupEnv("DOCKERIP")
	log.Println("val -> ", val)
	var cliErr error
	cli, cliErr = client.NewClient(val, "v1.30", nil, nil)
	if cliErr != nil {
		panic(cliErr)
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)

	}
	//createDockerBuild()
	containerId := "71e9127e0faf7e812f2892fa7d7ab95e6285344b0155a8b7155b706e659630ad"
	// arg0 := "soffice"
	// arg1 := "--invisible" //This command is optional, it will help to disable the splash screen of LibreOffice.
	// arg2 := "--convert-to"
	// arg3 := "pdf:writer_pdf_Export"
	// path := "/home/kiran/Downloads/keys/road.docx"
	// nout, err := exec.Command(arg0, arg1, arg2, arg3, path).Output()
	// if err != nil {
	// 	fmt.Println("nout err -> ", err)
	// }
	// fmt.Println(nout)

	// containerId, err := createNewContainer("codebox3")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("conatonerId Mani, ", containerId)

	// starterr := containerStart(containerId)
	// if starterr != nil {
	// 	panic(starterr)
	// }

	//StopRunningContainer(containerId)
	// err = copyPayloadFromHereToContainer(containerId, "/home/user/projects/", "../../../../../node/demo-docker2")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("err is not der, may be file has copied")

	//containerLogs(containerId)

	execerr := executeInsideContainer("main.js", "/home/user/projects/", "js", containerId)
	if execerr != nil {
		log.Panicln(err)
	}
	log.Println("check container")
}

func createNewContainer(image string) (string, error) {
	if image == "" {
		return "", errors.New("Image name cant be empty")
	}
	var timeout *int
	timeout = new(int)
	*timeout = 300

	container, err := cli.ContainerCreate(context.Background(), &container.Config{Image: image, NetworkDisabled: true, OpenStdin: true, Tty: true, StopTimeout: timeout, Shell: []string{"bin/bash", "-c", "touch appa.py"}}, &container.HostConfig{Resources: container.Resources{CPUShares: 200, Memory: 524288000, MemorySwap: 0}}, nil, "kiran3")
	if err != nil {
		panic(err)
	}
	log.Println("Created Container -< ", container)

	return container.ID, nil

}

func containerStart(containerId string) error {
	// check if contianer is already running and if it is stopped then only start the container
	err := cli.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})

	if err != nil {
		panic(err)

	}
	log.Println("container started if code is working")
	return nil
}

func createDockerBuild() {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "Dockerfile"
	dockerFileReader, err := os.Open("testing/Dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(ctx, dockerFileTarReader,
		types.ImageBuildOptions{Context: dockerFileTarReader, Dockerfile: "Dockerfile", Tags: []string{"codebox3"}, BuildArgs: map[string]string{"USER_ID": "4000", "GROUP_ID": "4000"}})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

func StopRunningContainer(containerId string) {
	containerInspectRes, err := cli.ContainerInspect(context.Background(), containerId)
	if err != nil {
		panic(err)
	}

	log.Println("containerInspectRes -> ", containerInspectRes.State.Running)

	// only if containerId is running then proceed to stop container
	if containerInspectRes.State.Running == false {
		log.Println("as per inspect method container is not running, ")
		return
	}

	containerStopErr := cli.ContainerStop(context.Background(), containerId, nil)
	if containerStopErr != nil {
		panic(containerStopErr)
	}

	log.Println("container has been stopped")

}

func copyPayloadFromHereToContainer(containerId, distname, localpath string) error {

	ctx := context.Background()

	archive, err := newTarArchiveFromPath(localpath)
	if err != nil {
		return err
	}

	err = cli.CopyToContainer(ctx, containerId, distname, archive, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
	if err != nil {
		return err
	}

	return nil
}

func newTarArchiveFromPath(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	})

	if ok != nil {
		return nil, ok
	}
	ok = tw.Close()
	if ok != nil {
		return nil, ok
	}
	return bufio.NewReader(&buf), nil
}

func containerLogs(containerID string) {
	// attachToContainerOptions := docker.AttachToContainerOptions{
	// 	Container:    containerID,
	// 	OutputStream: os.Stdout,
	// 	ErrorStream:  os.Stderr,
	// 	Logs:         true,
	// 	Stdout:       true,
	// 	Stderr:       true,
	// }

	val, err := cli.ContainerAttach(context.Background(), containerID, types.ContainerAttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
	if err != nil {
		log.Println("attachToContaienr err1 -> ", err)
	}
	log.Println("Val1 -> ", val)

}

func executeInsideContainer(filename, destination, language, containerID string) error {
	if language == "js" {
		filename = "main.js"
	}
	filename = "main.js"

	//postman json body for create exec in a container working to execute the bash script
	// {
	// 	"AttachStdin": true,
	// 	"AttachStdout": true,
	// 	"AttachStderr": true,
	// 	"DetachKeys": "ctrl-p,ctrl-q",
	// 	"Tty": true,
	// 	"Cmd": ["/bin/bash",
	// 	"/home/user/projects/compiler.sh"
	// 	],
	// 	"Env": [
	// 	"FOO=bar",
	// 	"BAZ=quux"
	// 	],
	// 	"Privileged":false,
	// 	"User":"user"
	// 	}

	//Cmd: []string{"bin/bash", "-c", "touch /home/user/projects/appa.txt"}
	// working -> Cmd: []string{"touch", "/home/user/projects/appa.txt"}
	//node main.js > main.log 2>&1
	createExec, err := cli.ContainerExecCreate(context.Background(), containerID, types.ExecConfig{AttachStdin: false, AttachStdout: true, AttachStderr: true, Detach: true, Tty: true, Cmd: []string{"/bin/bash", "/home/user/projects/compiler.sh"}, User: "user"})
	//node main.js > app.log 2>&1

	if err != nil {
		panic(err)
	}
	log.Println("createExec -> ", createExec)

	startExecErr := cli.ContainerExecStart(context.Background(), createExec.ID, types.ExecStartCheck{Detach: false, Tty: false})
	if startExecErr != nil {
		log.Panicln("startTExecErr -> ", startExecErr)
		return startExecErr
	}

	// inspect containerExec
	inspectResp, err := cli.ContainerExecInspect(context.Background(), createExec.ID)
	if err != nil {
		log.Panicln("inspectResp err -> ", err)
	}
	log.Println("inspectResp.ExitCode -> ", inspectResp.ExitCode)
	if inspectResp.ExitCode == 0 {
		//success call to function for copying log file to host server from container
		err := copyLogFileFromContainerToHostServer(filename, destination, language, containerID)
		if err != nil {
			log.Panicln("copyLogFileFromContainerToHostServer err -> ", err)
		}

	} else {
		log.Println(" inspectResp.ExitCode err -> ", inspectResp.ExitCode)
		log.Println(inspectResp.ExecID)
		return errors.New("Exit Code is not 0")
	}
	return nil
}

func copyLogFileFromContainerToHostServer(filename, destination, language, containerID string) error {

	respReader, pathStat, err := cli.CopyFromContainer(context.Background(), containerID, "/home/user/projects/app.log")
	if err != nil {
		log.Panicln("copyLogFileFromContainerToHostServer -> ", err)
	}

	defer respReader.Close()

	log.Println("pathStat -> ", pathStat)

	log.Println("respReader -> ", respReader)
	body, err := ioutil.ReadAll(respReader)
	if err != nil {
		log.Panicln("response read failed")
	}
	log.Println("body -> ", body)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(respReader)
	// log.Println("respReader -> ", buf.String())
	return nil
}
