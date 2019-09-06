package annotations

import (
	"fmt"
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"os"
	"code.cloudfoundry.org/cli/types"
)


type metadataFields struct {
	Labels map[string]types.NullString `json:"labels"`
	Annotations map[string]types.NullString `json:"annotations"`
}

type metadata struct {
	Metadata metadataFields  `json:"metadata"`
}

func newMetadata() metadata {
	return metadata{Metadata: metadataFields{ Labels: map[string]types.NullString{}, Annotations: map[string]types.NullString{}} }
}

func processArgs(args []string, keyCheck func(string, *metadata) error) (string, string, string, *metadata, error) {
	var resourceType string
	var resourceName string
	var stackOption string

	resourceType, args = args[0], args[1:]
	switch strings.ToLower(resourceType) {
	case "app":
	case "org":
	case "space":
	case "stack":
	case "buildpack":
		break
	default:
		return "", "", "", nil, fmt.Errorf("not a valid resource type: %s", resourceType)
	}

	resourceName, args = args[0], args[1:]

	m := newMetadata()

	for i, s := range (args) {
		if s == "-s" || s == "--stack" {
			if strings.ToLower(resourceType) != "buildpack" {
				return "", "", "", nil, fmt.Errorf("--stack not allowed with type %s", strings.ToLower(resourceType))
			}
			if i == len(args)-1 {
				return "", "", "", nil, fmt.Errorf("No stack for option '%s'", s)
			}
			stackOption = args[i+1]
			args[i+1] = ""
		} else if s[0] == '-' {
			return "", "", "", nil, fmt.Errorf("invalid option of '%s'", s)
		} else if s != "" {
			err := keyCheck(s, &m)
			if err != nil {
				return "", "", "", nil, err
			}
		}
	}
	return resourceType, resourceName, stackOption, &m, nil
}

func SetAnnotations(cliConnection plugin.CliConnection, args []string) error {
	resourceType, resourceName, stackOption, metadataPtr, err := processArgs(args, func(s string, metadataPtr *metadata) error {
		idx := strings.Index(s, "=")
		fmt.Fprintf(os.Stderr, "= at %d\n", idx)
		if idx == -1 {
			return fmt.Errorf("no value part given for annotation '%s'", s)
		} else if idx == 0 {
			return fmt.Errorf("no key part given for annotation '%s'", s)
		}
		metadataPtr.Metadata.Annotations[s[:idx]] = types.NewNullString(s[idx + 1:])
		return nil
	})
	if err != nil {
		return err
	}
	if len(metadataPtr.Metadata.Annotations) == 0 {
		return fmt.Errorf("no annotations specified")
	}

	guid, endpointName, err := getResource(cliConnection, resourceType, resourceName, stackOption)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(metadataPtr)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/v3/%s/%s", endpointName, guid)
	lines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", endpoint,
		"-X", "PATCH", "-d", string(payload))
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}

func UnsetAnnotations(cliConnection plugin.CliConnection, args []string) error {
	resourceType, resourceName, stackOption, metadataPtr, err := processArgs(args, func(s string, metadataPtr *metadata) error {
		idx := strings.Index(s, "=")
		fmt.Fprintf(os.Stderr, "= at %d\n", idx)
		if idx != -1 {
			return fmt.Errorf("annotation key '%s' cannot contain an '='", s)
		}
		metadataPtr.Metadata.Annotations[s] = types.NewNullString()
		return nil
	})
	if err != nil {
		return err
	}
	if len(metadataPtr.Metadata.Annotations) == 0 {
		return fmt.Errorf("No annotations specified")
	}

	guid, endpointName, err := getResource(cliConnection, resourceType, resourceName, stackOption)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(metadataPtr)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/v3/%s/%s", endpointName, guid)
	lines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", endpoint,
		"-X", "PATCH", "-d", string(payload))
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}

func ViewAnnotations(cliConnection plugin.CliConnection, args []string) error {
	resourceType, resourceName, stackOption, metadataPtr, err := processArgs(args, func(s string, metadataPtr *metadata) error {
		metadataPtr.Metadata.Annotations[s] = types.NewNullString()
		return nil
	})
	if err != nil {
		return err
	}
	if len(metadataPtr.Metadata.Annotations) != 0 {
		return fmt.Errorf("extra arguments specified")
	}

	guid, endpointName, err := getResource(cliConnection, resourceType, resourceName, stackOption)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("/v3/%s/%s", endpointName, guid)
	lines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", endpoint)
	if err != nil {
		return err
	}
	output := []byte(strings.Join(lines, "\n"))

	m := newMetadata()
	err = json.Unmarshal(output, &m)
	if err != nil {
		return err
	}

	fmt.Printf("%s: %s\n\n", resourceType, resourceName)
	fmt.Printf("Annotations:\n\n")
	for k, v := range m.Metadata.Annotations {
		fmt.Printf("%s: %s\n\n", k, v.Value)
	}
	return nil
}

func getResource(cliConnection plugin.CliConnection, resourceType, resourceName, stackOption string) (string, string, error) {
	var guid string
	var endpointName string
	switch strings.ToLower(resourceType) {
	case "app":
		obj, err := cliConnection.GetApp(resourceName)
		if err != nil {
			return "", "", err
		}
		guid = obj.Guid
	case "org":
		obj, err := cliConnection.GetOrg(resourceName)
		if err != nil {
			return "", "", err
		}
		guid = obj.Guid
		endpointName = "organizations"
	case "space":
		obj, err := cliConnection.GetSpace(resourceName)
		if err != nil {
			return "", "", err
		}
		guid = obj.Guid
	case "stack":
		lines, err := cliConnection.CliCommandWithoutTerminalOutput("stack", resourceName, "--guid")
		if err != nil {
			return "", "", err
		}
		guid = strings.TrimSpace(lines[0])
	case "buildpack":
		endpoint := "/v3/buildpacks?names=" + resourceName
		if stackOption != "" {
			endpoint += "&stacks=" + stackOption
		}
		lines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", endpoint)
		if err != nil {
			return "", "", err
		}

		var buildpacks struct {
			Resources []struct {
				Guid string 
			} `json:"resources"`
		}

		output := []byte(strings.Join(lines, "\n"))
		json.Unmarshal(output, &buildpacks)
		resources := buildpacks.Resources
		switch len(resources) {
		case 0:
			return "", "", fmt.Errorf("no buildpacks match %s (%s)", resourceName, stackOption)
		case 1:
			break
		default:
			return "", "", fmt.Errorf("too many buildpacks (%d) match %s", len(resources), resourceName)
		}
		guid = resources[0].Guid

	default:
		return "", "", fmt.Errorf("unrecognized type %s", resourceType)
	}
	if endpointName == "" {
		endpointName = strings.ToLower(resourceType) + "s"
	}
	return guid, endpointName, nil
}
