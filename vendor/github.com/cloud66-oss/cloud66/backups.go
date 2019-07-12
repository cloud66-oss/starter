package cloud66

import (
	"strconv"
	"strings"
	"time"
)

/*
 */
var BackupStatus = map[int]string{
	0: "Ok",     // BCK_OK
	1: "Failed", // BCK_FAILED
}

var RestoreStatus = map[int]string{
	0: "Not Restored",   // RST_NA
	1: "Restoring",      // RST_RESTORING
	2: "Restored",       // RST_OK
	3: "Restore Failed", // RST_FAILED
}

var VerifyStatus = map[int]string{
	0: "Not Verified",        // VRF_NA
	1: "Verifying",           // VRF_VERIFING
	2: "Verified",            // VRF_OK
	3: "Verification Failed", // VRF_FAILED
	4: "Unable to Verify",    // VRF_INTERNAL_ISSUE
}

type ManagedBackup struct {
	Id            int       `json:"id"`
	ServerUid     string    `json:"server_uid"`
	Filename      string    `json:"file_name"`
	DbType        string    `json:"db_type"`
	DatabaseName  string    `json:"database_name"`
	FileBase      string    `json:"file_base"`
	BackupDate    time.Time `json:"backup_date_iso"`
	BackupStatus  int       `json:"backup_status"`
	BackupResult  string    `json:"backup_result"`
	RestoreStatus int       `json:"restore_status"`
	RestoreResult string    `json:"restore_result"`
	CreatedAt     time.Time `json:"created_at_iso"`
	UpdatedAt     time.Time `json:"updated_at_iso"`
	VerifyStatus  int       `json:"verify_status"`
	VerifyResult  string    `json:"verify_result"`
	StoragePath   string    `json:"storage_path"`
	SkipTables    string    `json:"skip_tables"`
}

type BackupSegmentIndex struct {
	Filename  string `json:"name"`
	Extension string `json:"id"`
}

type BackupSegment struct {
	Ok  bool   `json:"ok"`
	Url string `json:"public_url"`
}

func (c *Client) GetBackupSegmentIndeces(stackUid string, backupId int) ([]BackupSegmentIndex, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []BackupSegmentIndex
	var backupSegIndex []BackupSegmentIndex

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/backups/"+strconv.Itoa(backupId)+"/files.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		backupSegIndex = nil
		err = c.DoReq(req, &backupSegIndex, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, backupSegIndex...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil

}

func (c *Client) GetBackupSegment(stackUid string, backupId int, extension string) (*BackupSegment, error) {
	ext := ""
	if extension != "" {
		ext = "/" + extension
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/backups/"+strconv.Itoa(backupId)+"/files/"+ext+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var backupSegmentRes *BackupSegment
	err = c.DoReq(req, &backupSegmentRes, nil)
	if err != nil {
		return nil, err
	}

	// fix percentage deserialize go bug
	backupSegmentRes.Url = strings.Replace(backupSegmentRes.Url, "%25", "%", -1)
	return backupSegmentRes, err

}

func (c *Client) NewBackup(stackUid string, dbtypes *string, frequency *string, keep *int, gzip *bool, exclude_tables *string, run_on_replica *bool, logical_backup *bool) error {

	params := struct {
		DbType        *string `json:"db_type"`
		Frequency     *string `json:"frequency"`
		KeepCount     *int    `json:"keep_count"`
		Gzip          *bool   `json:"gzip"`
		ExcludeTables *string `json:"excluded_tables"`
		RunOnReplica  *bool   `json:"run_on_replica_server"`
		LogicalBackup *bool   `json:"logical_backup"`
	}{
		DbType:        dbtypes,
		Frequency:     frequency,
		KeepCount:     keep,
		Gzip:          gzip,
		ExcludeTables: exclude_tables,
		RunOnReplica:  run_on_replica,
		LogicalBackup: logical_backup,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/backups.json", params, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil

}
