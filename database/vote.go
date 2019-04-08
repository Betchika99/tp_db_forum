package database

import "github.com/Betchika99/tp_db_project/models"

const (
	sqlInsertVote = `INSERT INTO votes (voice, nickname, thread_id)
					 VALUES ($1, $2, $3)
					 RETURNING voice, nickname`

	sqlSelectVote = `SELECT voice FROM votes WHERE nickname = $1 AND thread_id = $2`

	sqlUpdateVote = `UPDATE votes
					 SET voice = $1
					 WHERE nickname = $2 AND thread_id = $3
					 RETURNING voice, nickname`
)

func InsertVote(vote models.Vote, threadID int) (models.Vote, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return vote, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlInsertVote, vote.Voice, vote.Nickname, threadID)

	voteGot := models.Vote{}

	err = data.Scan(&voteGot.Voice, &voteGot.Nickname)
	if err != nil {
		return vote, err
	}

	if err = transaction.Commit(); err != nil {
		return vote, err
	}

	return voteGot, nil
}

func SelectVote(vote models.Vote, threadID int) (models.Vote, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return vote, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlSelectVote, vote.Nickname, threadID)

	voteGot := models.Vote{}
	err = data.Scan(&voteGot.Voice)
	if err != nil {
		return vote, err
	}

	if err = transaction.Commit(); err != nil {
		return vote, err
	}

	return voteGot, nil
}

func UpdateVote(vote models.Vote, threadID int) (models.Vote, error) {
	transaction, err := GetConnect().Begin()
	if err != nil {
		return vote, err
	}
	defer transaction.Rollback()

	data := transaction.QueryRow(sqlUpdateVote, vote.Voice, vote.Nickname, threadID)

	voteGot := models.Vote{}
	err = data.Scan(&voteGot.Voice, &voteGot.Nickname)
	if err != nil {
		return vote, err
	}

	if err = transaction.Commit(); err != nil {
		return vote, err
	}

	return voteGot, nil
}