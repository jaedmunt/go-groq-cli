# Strike
## A lightning fast AI Command line interface

> Note: The name "Strike" was chosen to avoid confusion with Groq's official releases and for its CLI-friendly nature.


#### ⚡ Speed Above All
- Built with [Go](https://go.dev/) for maximum performance and minimal footprint
- Powered by [Groq](https://groq.com/) for rapid, cost-effective responses
- Lightweight executable for instant access
</br>

### Quick Start
1. Build the executable:
    ```bash
    go build -o strike
    ```

2. Add to system PATH:
    ```bash
    # Windows
    move strike.exe %USERPROFILE%\go\bin

    # Linux/macOS
    mv strike $HOME/go/bin
    ```

3. Set Groq API key:
    ```bash
    # Windows
    setx GROQ_API_KEY "your-api-key"

    # Linux/macOS
    echo "export GROQ_API_KEY=your-api-key" >> ~/.bashrc
    source ~/.bashrc
    ```

4. Open command prompt (⊞ Win + R, type `CMD`)

5. Start chatting!
    ```bash
    strike "Tell me about Go programming"
    ```

Messages stream directly in your terminal.

