job "raw-exec-example" {
  type = "batch"
  group "raw" {
    task "raw" {
      driver = "raw_exec"
      config {
        command = "whoami"
      }
    }
  }
}