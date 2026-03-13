package envelope

// File transfer message types.
const (
	MessageTypeFileOffer        MessageType = "file_offer"
	MessageTypeFileAccept       MessageType = "file_accept"
	MessageTypeFileReject       MessageType = "file_reject"
	MessageTypeTransferReady    MessageType = "transfer_ready"
	MessageTypeTransferComplete MessageType = "transfer_complete"
	MessageTypeChunkAck         MessageType = "chunk_ack"
	MessageTypeResumeRequest    MessageType = "resume_request"
	MessageTypeFileChunk        MessageType = "file_chunk"
)
