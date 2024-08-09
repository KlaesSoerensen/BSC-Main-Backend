namespace BSC_Main_Backend.dto.response;

/// <param name="Code">Join Code</param>
/// <param name="LobbyId">Lobby id as known to the Multiplayer Broker</param>
/// <param name="MultiplayerBackendIp">Address of Multiplayer Broker</param>
public record OpenColonyResponseDTO(uint Code, uint LobbyId, string MultiplayerBackendIp);