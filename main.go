
package main

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
    "fmt"
)

type Slave struct {
	IP     string
	Active bool
}

var slaves = []Slave{
	{IP: "192.168.1.8"},
   // {IP: "192.168.80.61"},
	// Add more slaves if needed
}

var mu sync.Mutex

func main() {

	go monitorSlaves()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/shutdown", shutdownHandler)

	log.Println("Master is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func remoteShutdown(ip string) error {
	remoteMachine := fmt.Sprintf("\\\\%s", ip)
	cmd := exec.Command("shutdown", "/s", "/f", "/t", "0", "/m", remoteMachine)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Shutdown command failed: %s\n", string(output))
	}
	return err
}



func monitorSlaves() {
	for {
		mu.Lock()
		for i := range slaves {
			err := ping(slaves[i].IP)
			slaves[i].Active = (err == nil)
		}
		mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func ping(ip string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}
	return cmd.Run()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, slaves)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "IP address is missing", http.StatusBadRequest)
		return
	}

	log.Println("Shutdown requested for", ip)

	err := remoteShutdown(ip)
	if err != nil {
		log.Println("Failed to shutdown", ip, ":", err)
		http.Error(w, "Failed to shutdown slave", http.StatusInternalServerError)
		return
	}

	log.Println(ip, "shutdown command sent successfully.")
	w.WriteHeader(http.StatusOK)
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Master Control Panel</title>
    <style>
        body {
            background: linear-gradient(to right, #6a11cb,rgb(37, 252, 51));
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            color: white;
            text-align: center;
            padding: 50px;
        }
        table {
            width: 70%;
            margin: auto;
            border-collapse: collapse;
            margin-top: 30px;
            box-shadow: 0 8px 16px rgba(74, 72, 72, 0.3);
        }
        th, td {
            padding: 15px;
            background: rgba(247, 185, 185, 0.89);
            border-bottom: 1px solid rgba(255,255,255,0.3);
            font-size: 20px;
        }
        th {
            background: rgba(61, 150, 33, 0.3);
        }
        button {
            padding: 10px 20px;
            font-size: 16px;
            background-color: #ff416c;
            border: none;
            border-radius: 5px;
            color: white;
            cursor: pointer;
            transition: background 0.3s ease;
        }
        button:hover {
            background-color: #ff4b2b;
        }
        .active {
            color: #00ff00;
            font-weight: bold;
            animation: pulse 2s infinite;
        }
        .notactive {
            color: #ff0000;
            font-weight: bold;
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.1); }
            100% { transform: scale(1); }
        }
    </style>
</head>
<body>
    <h1>🌟<< Master Control Panel >>🌟</h1>
    <table>
        <thead>
            <tr>
                <th>IP Address</th>
                <th>Status</th>
                <th>Action</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.IP}}</td>
                <td class="{{if .Active}}active{{else}}notactive{{end}}">
                    {{if .Active}}Active{{else}}Not Active{{end}}
                </td>
                <td>
                    {{if .Active}}
                    <button onclick="shutdown('{{.IP}}')">Shutdown</button>
                    {{else}}
                    ---
                    {{end}}
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>

    <script>
        function shutdown(ip) {
            if(confirm("Are you sure you want to shutdown " + ip + "?")) {
                fetch("/shutdown?ip=" + ip)
                .then(response => {
                    if (response.ok) {
                        alert("Shutdown command sent!");
                        location.reload();
                    } else {
                        alert("Failed to shutdown!");
                    }
                })
                .catch(err => {
                    alert("Error: " + err);
                });
            }
        }
        setInterval(() => {
            location.reload();
        }, 10000);
    </script>
</body>
</html>
`
