package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func install_nodejs() {
	node_js_url := "https://nodejs.org/dist/v18.17.1/node-v18.17.1-x64.msi"
	download_res := downloadFile("node-setup.msi", node_js_url)
	if download_res == nil {
		println("Error: NodeJs Installation Failed")
	}

	pwd, _ := os.Getwd()
	cmd := exec.Command("msiexec", "/i", pwd+"\\node-setup.msi")
	_, exec_err := cmd.Output()
	if exec_err != nil {
		fmt.Println("An Error Occur: Fail to Install NodeJs ", exec_err.Error())
	}
}

func runCommandWithProgress(command string, args ...string) error {
	// Create a new command
	cmd := exec.Command(command, args...)

	// Redirect standard output and standard error to our program's streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func load_url(url string) {
	_, err := exec.Command("powershell", "-command", "start", url).Output()
	if err != nil {
		println("Fail to load URL", err.Error())
	}
}

func run_bot() {

	cdir_err := os.Chdir("./whatsapp-ai-bot-master")
	if cdir_err != nil {
		println("Fail to change directory '/whatsapp-ai-bot-master'")
	}

	// input keys

	var open_ai_api_key string
	var stability_ai_api_key string

	open_ai_url := "https://platform.openai.com/account/api-keys"
	stability_ai_url := "https://platform.stability.ai/docs/getting-started/authentication"

	fmt.Println("\nGet OpenAi API Key from " + open_ai_url)
	load_url(open_ai_url)

	print("Enter OpenAI API Key OR enter NONE to skip _ ")
	fmt.Scan(&open_ai_api_key)

	fmt.Println("\nGet StabilityAI API Key from " + stability_ai_url)
	exec.Command("start", open_ai_url)

	print("Enter StabilityAI API Key OR enter NONE to skip _ ")
	fmt.Scan(&stability_ai_api_key)

	// creates .env

	dot_env, dot_env_err := os.Create(".env")

	if dot_env_err != nil {
		println("Error: Fail to create .env file ", dot_env_err.Error())
	}

	dot_env.WriteString("OPENAI_API_KEY=" + open_ai_api_key + "\nDREAMSTUDIO_API_KEY=" + stability_ai_api_key)

	defer dot_env.Close()

	// install dependencies
	install_res := runCommandWithProgress("npx", "yarn")
	if install_res != nil {
		println("Error: fail to run 'npx yarn' ", install_res.Error())
	}

	// Setup Guide
	pwd, _ := os.Getwd()
	println("\n=== Whatsapp AI Bot Ready To Run ===")
	println("\r - To run next time go to 'whatsapp-ai-bot-master' folder & run setup.sh")
	println("\r OR ")
	println("\r - copy & paste following code in command prompt\n   cd " + pwd + " && npx yarn dev\n\n")

	// run bot
	run_res := runCommandWithProgress("npx", "yarn", "dev")
	if run_res != nil {
		println("Error: Fail to run 'npx yarn dev' ", run_res.Error())
	}
}

// driver code
func main() {
	cmd := exec.Command("node", "--version")
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("Node-Js not install in this machine. Going to install Node-Js \nPlease Wait it may take some time")
		install_nodejs()
	} else {
		fmt.Println("NodeJs Already Install")
	}

	// fetch repo
	rm_err := os.RemoveAll("./whatsapp-ai-bot-master")
	if rm_err != nil {
		println("Fail to remove previous Directories", string(rm_err.Error()))
	}

	fmt.Println("Downloading Github Repo ...")
	repo_url := "https://github.com/Zain-ul-din/whatsapp-ai-bot/archive/refs/heads/master.zip"
	download_res := downloadFile("whatsapp-ai-bot.zip", repo_url)
	if download_res != nil {
		println("Error: Fail to clone github repo from " + repo_url)
	}

	fmt.Println("Unzipping Source Code")
	unzip_cmd := exec.Command(
		"powershell",
		"-command",
		"Expand-Archive",
		"-Path",
		"./whatsapp-ai-bot.zip",
		"-DestinationPath",
		"./",
	)

	_, unzip_err := unzip_cmd.Output()
	if unzip_err != nil {
		println("Error: Fail to unzip file")
	}

	run_bot()

	// pause
	fmt.Scanln()

	defer os.Remove("./node-setup.msi")
	defer os.Remove("./whatsapp-ai-bot.zip")
}