namespace BSC_Main_Backend.dto;

/// <param name="Id">Id of this entry in the DB, of key and chosen value</param>
/// <param name="Key">Key of preference, e.g. "Language"</param>
/// <param name="ChosenValue">Value of key, e.g. "DE"</param>
/// <param name="AvailableValues">Joined on column from AvailablePreferences containing all possible values for that preference</param>
public record PlayerPreferenceDTO(uint Id, string Key, string ChosenValue, List<string> AvailableValues);

public record PlayerPreferencesResponseDTO(List<PlayerPreferenceDTO> Preferences);