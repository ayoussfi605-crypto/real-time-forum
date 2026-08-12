package db

import "database/sql"

type Reaction struct {
	ID        int
	PostID    int
	UserID    int
	Reaction  string
	CreatedAt string
}

type ReactionStats struct {
	Likes    int
	Dislikes int
}

func GetReaction(postID, userID int) (*Reaction, error) {
    var reaction Reaction

    err := DB.QueryRow(`
        SELECT id, post_id, user_id, reaction, created_at
        FROM post_reactions
        WHERE post_id = ? AND user_id = ?
    `,
        postID,
        userID,
    ).Scan(
        &reaction.ID,
        &reaction.PostID,
        &reaction.UserID,
        &reaction.Reaction,
        &reaction.CreatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }

        return nil, err
    }

    return &reaction, nil
}

func CreateReaction(postID, userID int, reactionType string) error {
	_, err := DB.Exec(`
		INSERT INTO post_reactions (post_id, user_id, reaction)
		VALUES (?, ?, ?)
	`, postID, userID, reactionType)

	return err
}

func UpdateReaction(postID, userID int, reactionType string) error {
	_, err := DB.Exec(`
		UPDATE post_reactions
		SET reaction = ?
		WHERE post_id = ? AND user_id = ?
	`, reactionType, postID, userID)

	return err
}

func DeleteReaction(postID, userID int) error {
	_, err := DB.Exec(`
		DELETE FROM post_reactions
		WHERE post_id = ? AND user_id = ?
	`, postID, userID)

	return err
}

func CountReactions(postID int) (*ReactionStats, error) {
	var stats ReactionStats

	err := DB.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN reaction = 'like' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN reaction = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM post_reactions
		WHERE post_id = ?
	`, postID).Scan(
		&stats.Likes,
		&stats.Dislikes,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
