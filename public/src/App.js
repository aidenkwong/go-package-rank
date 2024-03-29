import awesomeGoJson from "./awesome_go.json";
import { useState } from "react";

function App() {
  const keys = Object.keys(awesomeGoJson[0]);
  const [awesomeGo, setAwesomeGo] = useState(awesomeGoJson);
  const [sortBy, setSortBy] = useState({ key: keys[1], direction: "desc" });

  const max = Object.fromEntries(
    keys.map((key) => {
      return [key, Math.max(...awesomeGo.map((item) => item[key]))];
    })
  );

  const handleSort = (e) => {
    if (e.target.innerText === "Imported By") {
      if (sortBy.direction === "desc") {
        setAwesomeGo([
          ...awesomeGo.sort((a, b) => a.ImportedBy - b.ImportedBy),
        ]);
        setSortBy({ key: "ImportedBy", direction: "asc" });
      } else {
        setAwesomeGo([
          ...awesomeGo.sort((a, b) => b.ImportedBy - a.ImportedBy),
        ]);
        setSortBy({ key: "ImportedBy", direction: "desc" });
      }
    }
    if (e.target.innerText === "GitHub Stars") {
      if (sortBy.direction === "desc") {
        setAwesomeGo([
          ...awesomeGo.sort((a, b) => a.GitHubStar - b.GitHubStar),
        ]);
        setSortBy({ key: "GitHubStar", direction: "asc" });
      } else {
        setAwesomeGo([
          ...awesomeGo.sort((a, b) => b.GitHubStar - a.GitHubStar),
        ]);
        setSortBy({ key: "GitHubStar", direction: "desc" });
      }
    }
  };

  return (
    <div className="App" style={{ padding: 16 }}>
      <h1>Awesome Go Package Rank</h1>
      <p>
        This is a table of Go packages with number of imports and GitHub Stars,
        with reference to{" "}
        <a href="https://github.com/avelino/awesome-go">avelino/awesome-go</a>,
        which is a curated list of awesome Go frameworks, libraries and
        software.
      </p>
      <div>
        <table
          style={{
            width: "100%",
            border: "1px solid #dddddd",
          }}
        >
          <thead>
            <tr>
              <th>No.</th>
              {keys.map((key, i) => (
                <th
                  key={key}
                  style={{
                    textAlign: i === 0 ? "left" : "right",
                    cursor: i !== 0 ? "pointer" : "default",
                    border: "1px solid #dddddd",
                    padding: 4,
                  }}
                  onClick={handleSort}
                >
                  {key === "ImportedBy"
                    ? "Imported By"
                    : key === "GitHubStar"
                    ? "GitHub Stars"
                    : key}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {awesomeGo.map((row, idx) => (
              <tr key={row.Name}>
                <td style={{ textAlign: "center" }}>{idx + 1}</td>
                {keys.map((key, i) => (
                  <td
                    style={{
                      textAlign: i === 0 ? "left" : "right",
                      border: "1px solid #dddddd",
                      padding: 4,
                      backgroundColor:
                        key === "ImportedBy"
                          ? `hsl(120,${
                              (Math.log(row[key]) / Math.log(max[key])) * 100
                            }%,75%)`
                          : key === "GitHubStar"
                          ? `hsl(240,${
                              (Math.log(row[key]) / Math.log(max[key])) * 100
                            }%,75%)`
                          : "",
                    }}
                    key={key}
                  >
                    {key === "Name" ? (
                      <a href={"https://" + row[key]}>{row[key]}</a>
                    ) : (
                      row[key]
                    )}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
