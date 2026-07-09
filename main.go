package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	cli "github.com/urfave/cli/v3"
	"log"
	"os"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {

	cmd := &cli.Command{
		Name:    "slack plugin",
		Usage:   "slack plugin",
		Action:  run,
		Version: fmt.Sprintf("%s+%s", version, build),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "webhook",
				Usage:   "slack webhook url",
				Sources: cli.EnvVars("SLACK_WEBHOOK", "PLUGIN_WEBHOOK"),
			},
			&cli.StringFlag{
				Name:    "webhook.notice",
				Usage:   "slack webhook url for notices (PR failures and successes)",
				Sources: cli.EnvVars("PLUGIN_WEBHOOK_NOTICE", "SLACK_WEBHOOK_NOTICE"),
			},
			&cli.StringFlag{
				Name:    "webhook.alerts",
				Usage:   "slack webhook url for alerts (failures on the default branch)",
				Sources: cli.EnvVars("PLUGIN_WEBHOOK_ALERTS", "SLACK_WEBHOOK_ALERTS"),
			},
			&cli.StringFlag{
				Name:    "channel",
				Usage:   "slack channel",
				Sources: cli.EnvVars("PLUGIN_CHANNEL"),
			},
			&cli.StringFlag{
				Name:    "recipient",
				Usage:   "slack recipient",
				Sources: cli.EnvVars("PLUGIN_RECIPIENT"),
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "slack username",
				Sources: cli.EnvVars("PLUGIN_USERNAME"),
			},
			&cli.StringFlag{
				Name:    "template",
				Usage:   "slack template",
				Sources: cli.EnvVars("PLUGIN_TEMPLATE"),
			},
			&cli.StringFlag{
				Name:    "fallback",
				Usage:   "slack fallback",
				Sources: cli.EnvVars("PLUGIN_FALLBACK"),
			},
			&cli.BoolFlag{
				Name:    "link-names",
				Usage:   "slack link names",
				Sources: cli.EnvVars("PLUGIN_LINK_NAMES"),
			},
			&cli.StringFlag{
				Name:    "image",
				Usage:   "slack image url",
				Sources: cli.EnvVars("PLUGIN_IMAGE_URL"),
			},
			&cli.StringFlag{
				Name:    "color",
				Usage:   "slack color",
				Sources: cli.EnvVars("PLUGIN_COLOR"),
			},
			&cli.StringFlag{
				Name:    "icon.url",
				Usage:   "slack icon url",
				Sources: cli.EnvVars("PLUGIN_ICON_URL"),
			},
			&cli.StringFlag{
				Name:    "icon.emoji",
				Usage:   "slack emoji url",
				Sources: cli.EnvVars("PLUGIN_ICON_EMOJI"),
			},
			&cli.StringFlag{
				Name:    "repo.owner",
				Usage:   "repository owner",
				Sources: cli.EnvVars("CI_REPO_OWNER"),
			},
			&cli.StringFlag{
				Name:    "repo.name",
				Usage:   "repository name",
				Sources: cli.EnvVars("CI_REPO_NAME"),
			},
			&cli.StringFlag{
				Name:    "commit.sha",
				Usage:   "git commit sha",
				Sources: cli.EnvVars("CI_COMMIT_SHA"),
				Value:   "00000000",
			},
			&cli.StringFlag{
				Name:    "commit.ref",
				Value:   "refs/heads/master",
				Usage:   "git commit ref",
				Sources: cli.EnvVars("CI_COMMIT_REF"),
			},
			&cli.StringFlag{
				Name:    "commit.branch",
				Value:   "master",
				Usage:   "git commit branch",
				Sources: cli.EnvVars("CI_COMMIT_BRANCH"),
			},
			&cli.StringFlag{
				Name:    "commit.author",
				Usage:   "git author username",
				Sources: cli.EnvVars("CI_COMMIT_AUTHOR"),
			},
			&cli.StringFlag{
				Name:    "commit.author.email",
				Usage:   "git author email",
				Sources: cli.EnvVars("CI_COMMIT_AUTHOR_EMAIL"),
			},
			&cli.StringFlag{
				Name:    "commit.author.avatar",
				Usage:   "git author avatar",
				Sources: cli.EnvVars("CI_COMMIT_AUTHOR_AVATAR"),
			},
			&cli.StringFlag{
				Name:    "commit.author.name",
				Usage:   "git author name",
				Sources: cli.EnvVars("CI_COMMIT_AUTHOR_NAME"),
			},
			&cli.StringFlag{
				Name:    "commit.pull",
				Usage:   "git pull request",
				Sources: cli.EnvVars("CI_COMMIT_PULL_REQUEST"),
			},
			&cli.StringFlag{
				Name:    "commit.pull.branch",
				Usage:   "git pull request branch",
				Sources: cli.EnvVars("CI_COMMIT_SOURCE_BRANCH"),
			},
			&cli.StringFlag{
				Name:    "commit.message",
				Usage:   "commit message",
				Sources: cli.EnvVars("CI_COMMIT_MESSAGE"),
			},
			&cli.StringFlag{
				Name:    "build.event",
				Value:   "push",
				Usage:   "build event",
				Sources: cli.EnvVars("CI_PIPELINE_EVENT"),
			},
			&cli.IntFlag{
				Name:    "build.number",
				Usage:   "build number",
				Sources: cli.EnvVars("CI_PIPELINE_NUMBER"),
			},
			&cli.IntFlag{
				Name:    "build.parent",
				Usage:   "build parent",
				Sources: cli.EnvVars("CI_PIPELINE_PARENT"),
			},
			&cli.StringFlag{
				Name:    "build.status",
				Usage:   "build status",
				Value:   "success",
				Sources: cli.EnvVars("PLUGIN_STATUS"),
			},
			&cli.StringFlag{
				Name:    "build.link",
				Usage:   "build link",
				Sources: cli.EnvVars("CI_PIPELINE_URL"),
			},
			&cli.Int64Flag{
				Name:    "build.started",
				Usage:   "build started",
				Sources: cli.EnvVars("CI_PIPELINE_STARTED"),
			},
			&cli.Int64Flag{
				Name:    "build.created",
				Usage:   "build created",
				Sources: cli.EnvVars("CI_PIPELINE_CREATED"),
			},
			&cli.StringFlag{
				Name:    "build.tag",
				Usage:   "build tag",
				Sources: cli.EnvVars("CI_COMMIT_TAG"),
			},
			&cli.StringFlag{
				Name:    "build.deployTo",
				Usage:   "environment deployed to",
				Sources: cli.EnvVars("CI_PIPELINE_DEPLOY_TARGET"),
			},
			&cli.Int64Flag{
				Name:    "job.started",
				Usage:   "job started",
				Sources: cli.EnvVars("CI_PIPELINE_STARTED"),
			},
			&cli.StringFlag{
				Name:    "custom.block",
				Usage:   "custom block to send to slack. ",
				Sources: cli.EnvVars("PLUGIN_CUSTOM_BLOCK"),
			},
			&cli.StringFlag{
				Name:    "access.token",
				Usage:   "slack access token",
				Sources: cli.EnvVars("PLUGIN_ACCESS_TOKEN", "SLACK_ACCESS_TOKEN"),
			},
			&cli.StringFlag{
				Name:    "mentions",
				Usage:   "slack mentions for the message.",
				Sources: cli.EnvVars("PLUGIN_MENTIONS"),
			},
			&cli.StringFlag{
				Name:    "custom.template",
				Usage:   "prebuilt custom template for the message.",
				Sources: cli.EnvVars("PLUGIN_CUSTOM_TEMPLATE"),
			},
			&cli.StringFlag{
				Name:    "message",
				Usage:   "slack message. either this or the custom template must be set. ",
				Sources: cli.EnvVars("PLUGIN_MESSAGE"),
			},

			// File send params
			&cli.StringFlag{
				Name:    "filepath",
				Usage:   "slack file path",
				Sources: cli.EnvVars("PLUGIN_FILE_PATH"),
			},
			&cli.StringFlag{
				Name:    "filename",
				Usage:   "slack file name",
				Sources: cli.EnvVars("PLUGIN_FILE_NAME"),
			},
			&cli.StringFlag{
				Name:    "title",
				Usage:   "slack title",
				Sources: cli.EnvVars("PLUGIN_TITLE"),
			},
			&cli.StringFlag{
				Name:    "initial_comment",
				Usage:   "slack initial comment",
				Sources: cli.EnvVars("PLUGIN_INITIAL_COMMENT"),
			},
			&cli.BoolFlag{
				Name:    "fail_on_error",
				Usage:   "fail build on error",
				Sources: cli.EnvVars("PLUGIN_FAIL_ON_ERROR"),
			},
			&cli.StringFlag{
				Name:    "slack_id_of",
				Usage:   "slack id required for the user email id",
				Sources: cli.EnvVars("PLUGIN_SLACK_USER_EMAIL_ID"),
			},
			&cli.StringFlag{
				Name:    "committer_list_git_path",
				Usage:   "git repo path holding the committers email id to fetch slack IDs from",
				Sources: cli.EnvVars("PLUGIN_GIT_REPO_PATH"),
			},
			&cli.BoolFlag{
				Name:    "plugin_committer_slack_id",
				Usage:   "flag to enable fetching slack IDs from the committers list",
				Sources: cli.EnvVars("PLUGIN_COMMITTERS_SLACK_ID"),
			},
			&cli.StringFlag{
				Name:    "description",
				Usage:   "description of the build",
				Sources: cli.EnvVars("PLUGIN_DESCRIPTION"),
			},
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(_ context.Context, cmd *cli.Command) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: cmd.String("repo.owner"),
			Name:  cmd.String("repo.name"),
		},
		Build: Build{
			Tag:    cmd.String("build.tag"),
			Number: cmd.Int("build.number"),
			Parent: cmd.Int("build.parent"),
			Event:  cmd.String("build.event"),
			Status: cmd.String("build.status"),
			Commit: cmd.String("commit.sha"),
			Ref:    cmd.String("commit.ref"),
			Branch: cmd.String("commit.branch"),
			Author: Author{
				Username: cmd.String("commit.author"),
				Name:     cmd.String("commit.author.name"),
				Email:    cmd.String("commit.author.email"),
				Avatar:   cmd.String("commit.author.avatar"),
			},
			Pull:       cmd.String("commit.pull"),
			PullBranch: cmd.String("commit.pull.branch"),
			Message:    newCommitMessage(cmd.String("commit.message")),
			DeployTo:   cmd.String("build.deployTo"),
			Link:       cmd.String("build.link"),
			Started:    cmd.Int64("build.started"),
			Created:    cmd.Int64("build.created"),
		},
		Job: Job{
			Started: cmd.Int64("job.started"),
		},
		Config: Config{
			Webhook:        cmd.String("webhook"),
			WebhookNotice:  cmd.String("webhook.notice"),
			WebhookAlerts:  cmd.String("webhook.alerts"),
			Channel:        cmd.String("channel"),
			Recipient:      cmd.String("recipient"),
			Username:       cmd.String("username"),
			Template:       cmd.String("template"),
			Fallback:       cmd.String("fallback"),
			ImageURL:       cmd.String("image"),
			IconURL:        cmd.String("icon.url"),
			IconEmoji:      cmd.String("icon.emoji"),
			Color:          cmd.String("color"),
			LinkNames:      cmd.Bool("link-names"),
			CustomBlock:    cmd.String("custom.block"),
			AccessToken:    cmd.String("access.token"),
			Mentions:       cmd.String("mentions"),
			CustomTemplate: cmd.String("custom.template"),
			Message:        cmd.String("message"),
			// File upload attributes
			FilePath:             cmd.String("filepath"),
			FileName:             cmd.String("filename"),
			Title:                cmd.String("title"),
			InitialComment:       cmd.String("initial_comment"),
			FailOnError:          cmd.Bool("fail_on_error"),
			SlackIdOf:            cmd.String("slack_id_of"),
			CommitterListGitPath: cmd.String("committer_list_git_path"),
			CommitterSlackId:     cmd.Bool("plugin_committer_slack_id"),
			Description:          cmd.String("description"),
		},
	}

	if plugin.Build.Commit == "" {
		plugin.Build.Commit = "0000000000000000000000000000000000000000"
	}
	if plugin.Config.Webhook == "" && plugin.Config.WebhookNotice == "" && plugin.Config.WebhookAlerts == "" && plugin.Config.AccessToken == "" {
		return errors.New("you must provide a webhook url or access token")
	}

	return plugin.Exec()
}
